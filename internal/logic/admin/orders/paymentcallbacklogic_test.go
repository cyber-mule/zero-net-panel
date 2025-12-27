package orders

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/bootstrap/migrations"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func setupPaymentCallbackTest(t *testing.T) (*svc.ServiceContext, func()) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	_, err = migrations.Apply(context.Background(), db, 0, false)
	require.NoError(t, err)

	repos, err := repository.NewRepositories(db)
	require.NoError(t, err)

	svcCtx := &svc.ServiceContext{
		DB:           db,
		Repositories: repos,
	}

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	return svcCtx, cleanup
}

func seedDefaultTemplate(t *testing.T, db *gorm.DB) repository.SubscriptionTemplate {
	t.Helper()

	now := time.Now().UTC()
	tpl := repository.SubscriptionTemplate{
		Name:        "Default Template",
		Description: "Test template",
		ClientType:  "clash",
		Format:      "go_template",
		Content:     "test",
		IsDefault:   true,
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: &now,
	}

	var existing repository.SubscriptionTemplate
	err := db.Where("name = ? AND client_type = ?", tpl.Name, tpl.ClientType).First(&existing).Error
	if err == nil {
		return existing
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		require.NoError(t, err)
	}

	require.NoError(t, db.Create(&tpl).Error)
	return tpl
}

func TestPaymentCallbackLogic_Success(t *testing.T) {
	svcCtx, cleanup := setupPaymentCallbackTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	customer := repository.User{
		Email:       "customer@test.dev",
		DisplayName: "Customer",
		Roles:       []string{"user"},
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, svcCtx.DB.Create(&customer).Error)

	plan := repository.Plan{
		Name:              "Premium",
		Slug:              "premium",
		Description:       "Premium plan",
		PriceCents:        3200,
		Currency:          "CNY",
		DurationDays:      30,
		TrafficLimitBytes: 4096,
		DevicesLimit:      5,
		Status:            "active",
		Visible:           true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	require.NoError(t, svcCtx.DB.Create(&plan).Error)
	seedDefaultTemplate(t, svcCtx.DB)

	orderRepo := svcCtx.Repositories.Order

	order, items, err := orderRepo.Create(ctx, repository.Order{
		UserID:        customer.ID,
		PlanID:        &plan.ID,
		Status:        repository.OrderStatusPendingPayment,
		PaymentMethod: repository.PaymentMethodExternal,
		PaymentStatus: repository.OrderPaymentStatusPending,
		TotalCents:    plan.PriceCents,
		Currency:      plan.Currency,
		PlanSnapshot: map[string]any{
			"name":                plan.Name,
			"duration_days":       plan.DurationDays,
			"traffic_limit_bytes": plan.TrafficLimitBytes,
			"devices_limit":       plan.DevicesLimit,
		},
	}, []repository.OrderItem{{
		ItemType:       "plan",
		ItemID:         plan.ID,
		Name:           plan.Name,
		Quantity:       1,
		UnitPriceCents: plan.PriceCents,
		Currency:       plan.Currency,
		SubtotalCents:  plan.PriceCents,
		Metadata: map[string]any{
			"duration_days":       plan.DurationDays,
			"traffic_limit_bytes": plan.TrafficLimitBytes,
			"devices_limit":       plan.DevicesLimit,
		},
		CreatedAt:      now,
	}})
	require.NoError(t, err)
	require.Len(t, items, 1)

	payment, err := orderRepo.CreatePayment(ctx, repository.OrderPayment{
		OrderID:     order.ID,
		Provider:    "stripe",
		Method:      repository.PaymentMethodExternal,
		Status:      repository.OrderPaymentStatusPending,
		AmountCents: plan.PriceCents,
		Currency:    plan.Currency,
	})
	require.NoError(t, err)

	logic := NewPaymentCallbackLogic(ctx, svcCtx)
	paidAt := now.Add(2 * time.Minute).Unix()
	resp, err := logic.Process(&types.AdminPaymentCallbackRequest{
		OrderID:   order.ID,
		PaymentID: payment.ID,
		Status:    repository.OrderPaymentStatusSucceeded,
		Reference: "gateway-ref",
		PaidAt:    &paidAt,
	})
	require.NoError(t, err)

	require.Equal(t, repository.OrderStatusPaid, resp.Order.OrderDetail.Status)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, resp.Order.OrderDetail.PaymentStatus)
	require.Equal(t, "gateway-ref", resp.Order.OrderDetail.PaymentReference)
	require.NotNil(t, resp.Order.OrderDetail.PaidAt)
	require.Len(t, resp.Order.OrderDetail.Payments, 1)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, resp.Order.OrderDetail.Payments[0].Status)

	storedOrder, _, err := orderRepo.Get(ctx, order.ID)
	require.NoError(t, err)
	require.Equal(t, repository.OrderStatusPaid, storedOrder.Status)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, storedOrder.PaymentStatus)
	require.Equal(t, "gateway-ref", storedOrder.PaymentReference)

	subs, _, err := svcCtx.Repositories.Subscription.ListByUser(ctx, customer.ID, repository.ListSubscriptionsOptions{})
	require.NoError(t, err)
	require.Len(t, subs, 1)
	require.Equal(t, plan.Name, subs[0].PlanName)

	paymentsMap, err := orderRepo.ListPayments(ctx, []uint64{order.ID})
	require.NoError(t, err)
	require.Len(t, paymentsMap[order.ID], 1)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, paymentsMap[order.ID][0].Status)
}

func TestPaymentCallbackLogic_Failed(t *testing.T) {
	svcCtx, cleanup := setupPaymentCallbackTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	customer := repository.User{
		Email:       "customer2@test.dev",
		DisplayName: "Customer 2",
		Roles:       []string{"user"},
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, svcCtx.DB.Create(&customer).Error)

	orderRepo := svcCtx.Repositories.Order

	order, _, err := orderRepo.Create(ctx, repository.Order{
		UserID:        customer.ID,
		Status:        repository.OrderStatusPendingPayment,
		PaymentMethod: repository.PaymentMethodExternal,
		PaymentStatus: repository.OrderPaymentStatusPending,
		TotalCents:    4200,
		Currency:      "CNY",
	}, []repository.OrderItem{})
	require.NoError(t, err)

	payment, err := orderRepo.CreatePayment(ctx, repository.OrderPayment{
		OrderID:     order.ID,
		Provider:    "alipay",
		Method:      repository.PaymentMethodExternal,
		Status:      repository.OrderPaymentStatusPending,
		AmountCents: 4200,
		Currency:    "CNY",
	})
	require.NoError(t, err)

	logic := NewPaymentCallbackLogic(ctx, svcCtx)
	failureCode := "timeout"
	failureMsg := "payment timeout"
	resp, err := logic.Process(&types.AdminPaymentCallbackRequest{
		OrderID:        order.ID,
		PaymentID:      payment.ID,
		Status:         repository.OrderPaymentStatusFailed,
		FailureCode:    failureCode,
		FailureMessage: failureMsg,
	})
	require.NoError(t, err)

	require.Equal(t, repository.OrderStatusPaymentFailed, resp.Order.OrderDetail.Status)
	require.Equal(t, repository.OrderPaymentStatusFailed, resp.Order.OrderDetail.PaymentStatus)
	require.Equal(t, failureCode, resp.Order.OrderDetail.PaymentFailureCode)
	require.Equal(t, failureMsg, resp.Order.OrderDetail.PaymentFailureMessage)
	require.Len(t, resp.Order.OrderDetail.Payments, 1)
	require.Equal(t, repository.OrderPaymentStatusFailed, resp.Order.OrderDetail.Payments[0].Status)

	storedOrder, _, err := orderRepo.Get(ctx, order.ID)
	require.NoError(t, err)
	require.Equal(t, repository.OrderStatusPaymentFailed, storedOrder.Status)
	require.Equal(t, repository.OrderPaymentStatusFailed, storedOrder.PaymentStatus)
	require.Equal(t, failureCode, storedOrder.PaymentFailureCode)
	require.Equal(t, failureMsg, storedOrder.PaymentFailureReason)

	paymentsMap, err := orderRepo.ListPayments(ctx, []uint64{order.ID})
	require.NoError(t, err)
	require.Equal(t, repository.OrderPaymentStatusFailed, paymentsMap[order.ID][0].Status)
}

func TestPaymentCallbackLogic_Idempotent(t *testing.T) {
	svcCtx, cleanup := setupPaymentCallbackTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	customer := repository.User{
		Email:       "customer3@test.dev",
		DisplayName: "Customer 3",
		Roles:       []string{"user"},
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, svcCtx.DB.Create(&customer).Error)

	plan := repository.Plan{
		Name:              "Basic",
		Slug:              "basic",
		Description:       "Basic plan",
		PriceCents:        1000,
		Currency:          "CNY",
		DurationDays:      30,
		TrafficLimitBytes: 1024,
		DevicesLimit:      2,
		Status:            "active",
		Visible:           true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	require.NoError(t, svcCtx.DB.Create(&plan).Error)
	seedDefaultTemplate(t, svcCtx.DB)

	orderRepo := svcCtx.Repositories.Order

	order, _, err := orderRepo.Create(ctx, repository.Order{
		UserID:        customer.ID,
		PlanID:        &plan.ID,
		Status:        repository.OrderStatusPendingPayment,
		PaymentMethod: repository.PaymentMethodExternal,
		PaymentStatus: repository.OrderPaymentStatusPending,
		TotalCents:    plan.PriceCents,
		Currency:      plan.Currency,
		PlanSnapshot: map[string]any{
			"name":                plan.Name,
			"duration_days":       plan.DurationDays,
			"traffic_limit_bytes": plan.TrafficLimitBytes,
			"devices_limit":       plan.DevicesLimit,
		},
	}, []repository.OrderItem{{
		ItemType:       "plan",
		ItemID:         plan.ID,
		Name:           plan.Name,
		Quantity:       1,
		UnitPriceCents: plan.PriceCents,
		Currency:       plan.Currency,
		SubtotalCents:  plan.PriceCents,
		Metadata: map[string]any{
			"duration_days":       plan.DurationDays,
			"traffic_limit_bytes": plan.TrafficLimitBytes,
			"devices_limit":       plan.DevicesLimit,
		},
	}})
	require.NoError(t, err)

	payment, err := orderRepo.CreatePayment(ctx, repository.OrderPayment{
		OrderID:     order.ID,
		Provider:    "alipay",
		Method:      repository.PaymentMethodExternal,
		Status:      repository.OrderPaymentStatusPending,
		AmountCents: plan.PriceCents,
		Currency:    plan.Currency,
	})
	require.NoError(t, err)

	logic := NewPaymentCallbackLogic(ctx, svcCtx)
	resp, err := logic.Process(&types.AdminPaymentCallbackRequest{
		OrderID:        order.ID,
		PaymentID:      payment.ID,
		Status:         repository.OrderPaymentStatusSucceeded,
		Reference:      "ref-1",
		FailureCode:    "",
		FailureMessage: "",
	})
	require.NoError(t, err)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, resp.Order.OrderDetail.PaymentStatus)

	// Same status should be treated idempotently and return existing record.
	second, err := logic.Process(&types.AdminPaymentCallbackRequest{
		OrderID:        order.ID,
		PaymentID:      payment.ID,
		Status:         repository.OrderPaymentStatusSucceeded,
		Reference:      "ref-1",
		FailureCode:    "",
		FailureMessage: "",
	})
	require.NoError(t, err)
	require.Equal(t, resp.Order.OrderDetail.PaymentStatus, second.Order.OrderDetail.PaymentStatus)
	require.Equal(t, resp.Order.OrderDetail.PaymentReference, second.Order.OrderDetail.PaymentReference)

	// Downgrade from success to failed should be rejected.
	_, err = logic.Process(&types.AdminPaymentCallbackRequest{
		OrderID:        order.ID,
		PaymentID:      payment.ID,
		Status:         repository.OrderPaymentStatusFailed,
		FailureCode:    "duplicate",
		FailureMessage: "should not override success",
	})
	require.Error(t, err)
}
