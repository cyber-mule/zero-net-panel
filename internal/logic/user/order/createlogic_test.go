package order

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/bootstrap/migrations"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/testutil"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func setupCreateLogicTest(t *testing.T) (*svc.ServiceContext, func()) {
	t.Helper()

	testutil.RequireSQLite(t)

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

func seedPaymentChannel(t *testing.T, db *gorm.DB, code string, enabled bool, config map[string]any) repository.PaymentChannel {
	t.Helper()

	now := time.Now().UTC()
	channel := repository.PaymentChannel{
		Name:      code,
		Code:      code,
		Provider:  code,
		Enabled:   enabled,
		SortOrder: 1,
		Config:    config,
		CreatedAt: now,
		UpdatedAt: now,
	}
	require.NoError(t, db.Create(&channel).Error)
	return channel
}

func TestCreateOrderWithBalancePayment(t *testing.T) {
	svcCtx, cleanup := setupCreateLogicTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	user := repository.User{
		Email:       "buyer@test.dev",
		DisplayName: "Buyer",
		Roles:       []string{"user"},
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, svcCtx.DB.Create(&user).Error)

	plan := repository.Plan{
		Name:              "Standard",
		Slug:              "standard",
		Description:       "Standard plan",
		PriceCents:        1500,
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
	template := seedDefaultTemplate(t, svcCtx.DB)

	balanceRepo := svcCtx.Repositories.Balance
	_, _, err := balanceRepo.ApplyTransaction(ctx, user.ID, repository.BalanceTransaction{
		Type:        "recharge",
		AmountCents: 5000,
		Currency:    "CNY",
		Reference:   "seed",
		Description: "seed balance",
	})
	require.NoError(t, err)

	claims := security.UserClaims{ID: user.ID, Email: user.Email, Roles: []string{"user"}}
	reqCtx := security.WithUser(ctx, claims)

	logic := NewCreateLogic(reqCtx, svcCtx)
	resp, err := logic.Create(&types.UserCreateOrderRequest{
		PlanID: plan.ID,
	})
	require.NoError(t, err)

	require.Equal(t, repository.OrderStatusPaid, resp.Order.Status)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, resp.Order.PaymentStatus)
	require.Equal(t, repository.PaymentMethodBalance, resp.Order.PaymentMethod)
	require.NotNil(t, resp.Order.PaidAt)
	require.NotNil(t, resp.Transaction)
	require.Equal(t, int64(-1500), resp.Transaction.AmountCents)
	require.Empty(t, resp.Order.Payments)

	balanceSnapshot := resp.Balance
	require.Equal(t, user.ID, balanceSnapshot.UserID)
	require.Equal(t, int64(3500), balanceSnapshot.BalanceCents)

	storedOrder, _, err := svcCtx.Repositories.Order.Get(ctx, resp.Order.ID)
	require.NoError(t, err)
	require.Equal(t, repository.OrderStatusPaid, storedOrder.Status)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, storedOrder.PaymentStatus)

	subs, _, err := svcCtx.Repositories.Subscription.ListByUser(ctx, user.ID, repository.ListSubscriptionsOptions{})
	require.NoError(t, err)
	require.Len(t, subs, 1)
	require.Equal(t, plan.Name, subs[0].PlanName)
	require.Equal(t, plan.DevicesLimit, subs[0].DevicesLimit)
	require.Equal(t, plan.TrafficLimitBytes, subs[0].TrafficTotalBytes)
	require.Equal(t, template.ID, subs[0].TemplateID)
}

func TestCreateOrderWithExternalPayment(t *testing.T) {
	svcCtx, cleanup := setupCreateLogicTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.NotEmpty(t, payload["order_no"])
		require.NotEmpty(t, payload["amount"])
		require.NotEmpty(t, payload["notify_url"])
		require.NotEmpty(t, payload["return_url"])
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"data":{"pay_url":"https://pay.test/redirect","reference":"ref-001"}}`))
		require.NoError(t, err)
	}))
	defer gateway.Close()

	user := repository.User{
		Email:       "buyer2@test.dev",
		DisplayName: "Buyer 2",
		Roles:       []string{"user"},
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, svcCtx.DB.Create(&user).Error)

	plan := repository.Plan{
		Name:              "Premium",
		Slug:              "premium",
		Description:       "Premium plan",
		PriceCents:        2600,
		Currency:          "CNY",
		DurationDays:      30,
		TrafficLimitBytes: 2048,
		DevicesLimit:      3,
		Status:            "active",
		Visible:           true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	require.NoError(t, svcCtx.DB.Create(&plan).Error)
	seedDefaultTemplate(t, svcCtx.DB)
	seedPaymentChannel(t, svcCtx.DB, "stripe", true, map[string]any{
		"mode":       "http",
		"notify_url": "https://notify.test/callback?order_id={{order_id}}&payment_id={{payment_id}}",
		"return_url": "https://return.test/orders/{{order_number}}",
		"http": map[string]any{
			"endpoint":  gateway.URL,
			"method":    "POST",
			"body_type": "json",
			"payload": map[string]any{
				"order_no":   "{{order_number}}",
				"amount":     "{{amount}}",
				"notify_url": "{{notify_url}}",
				"return_url": "{{return_url}}",
			},
		},
		"response": map[string]any{
			"pay_url":   "data.pay_url",
			"reference": "data.reference",
		},
	})

	claims := security.UserClaims{ID: user.ID, Email: user.Email, Roles: []string{"user"}}
	reqCtx := security.WithUser(ctx, claims)

	logic := NewCreateLogic(reqCtx, svcCtx)
	resp, err := logic.Create(&types.UserCreateOrderRequest{
		PlanID:           plan.ID,
		PaymentMethod:    repository.PaymentMethodExternal,
		PaymentChannel:   "stripe",
		PaymentReturnURL: "https://example.com/return",
	})
	require.NoError(t, err)

	require.Equal(t, repository.OrderStatusPendingPayment, resp.Order.Status)
	require.Equal(t, repository.OrderPaymentStatusPending, resp.Order.PaymentStatus)
	require.Equal(t, repository.PaymentMethodExternal, resp.Order.PaymentMethod)
	require.NotEmpty(t, resp.Order.PaymentIntentID)
	require.Nil(t, resp.Transaction)
	require.Equal(t, int64(0), resp.Balance.BalanceCents)
	require.Len(t, resp.Order.Payments, 1)

	payment := resp.Order.Payments[0]
	require.Equal(t, repository.OrderPaymentStatusPending, payment.Status)
	require.Equal(t, plan.PriceCents, payment.AmountCents)
	require.Equal(t, plan.Currency, payment.Currency)
	require.Equal(t, "stripe", payment.Provider)
	require.Equal(t, "https://pay.test/redirect", payment.Metadata["pay_url"])

	storedOrder, _, err := svcCtx.Repositories.Order.Get(ctx, resp.Order.ID)
	require.NoError(t, err)
	require.Equal(t, repository.OrderStatusPendingPayment, storedOrder.Status)
	require.Equal(t, repository.OrderPaymentStatusPending, storedOrder.PaymentStatus)
	require.Equal(t, repository.PaymentMethodExternal, storedOrder.PaymentMethod)

	paymentsMap, err := svcCtx.Repositories.Order.ListPayments(ctx, []uint64{storedOrder.ID})
	require.NoError(t, err)
	require.Len(t, paymentsMap[storedOrder.ID], 1)
	require.Equal(t, repository.OrderPaymentStatusPending, paymentsMap[storedOrder.ID][0].Status)
}

func TestCreateOrderWithExternalPaymentInvalidChannel(t *testing.T) {
	cases := []struct {
		name string
		seed func(t *testing.T, db *gorm.DB)
		code string
	}{
		{
			name: "missing",
			code: "stripe",
		},
		{
			name: "disabled",
			code: "stripe",
			seed: func(t *testing.T, db *gorm.DB) {
				seedPaymentChannel(t, db, "stripe", false, nil)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svcCtx, cleanup := setupCreateLogicTest(t)
			defer cleanup()

			ctx := context.Background()
			now := time.Now().UTC()

			user := repository.User{
				Email:       "buyer-invalid@test.dev",
				DisplayName: "Buyer Invalid",
				Roles:       []string{"user"},
				Status:      "active",
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			require.NoError(t, svcCtx.DB.Create(&user).Error)

			plan := repository.Plan{
				Name:              "Invalid Plan",
				Slug:              "invalid-plan",
				Description:       "Plan",
				PriceCents:        1200,
				Currency:          "CNY",
				DurationDays:      30,
				TrafficLimitBytes: 1024,
				DevicesLimit:      1,
				Status:            "active",
				Visible:           true,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			require.NoError(t, svcCtx.DB.Create(&plan).Error)
			if tc.seed != nil {
				tc.seed(t, svcCtx.DB)
			}

			claims := security.UserClaims{ID: user.ID, Email: user.Email, Roles: []string{"user"}}
			reqCtx := security.WithUser(ctx, claims)

			logic := NewCreateLogic(reqCtx, svcCtx)
			resp, err := logic.Create(&types.UserCreateOrderRequest{
				PlanID:         plan.ID,
				PaymentMethod:  repository.PaymentMethodExternal,
				PaymentChannel: tc.code,
			})
			require.Error(t, err)
			require.ErrorIs(t, err, repository.ErrInvalidArgument)
			require.Nil(t, resp)
		})
	}
}

func TestCreateOrderWithManualPayment(t *testing.T) {
	svcCtx, cleanup := setupCreateLogicTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	user := repository.User{
		Email:       "buyer-manual@test.dev",
		DisplayName: "Buyer Manual",
		Roles:       []string{"user"},
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, svcCtx.DB.Create(&user).Error)

	plan := repository.Plan{
		Name:              "Manual Plan",
		Slug:              "manual-plan",
		Description:       "Manual payment plan",
		PriceCents:        1800,
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

	claims := security.UserClaims{ID: user.ID, Email: user.Email, Roles: []string{"user"}}
	reqCtx := security.WithUser(ctx, claims)

	logic := NewCreateLogic(reqCtx, svcCtx)
	resp, err := logic.Create(&types.UserCreateOrderRequest{
		PlanID:        plan.ID,
		PaymentMethod: repository.PaymentMethodManual,
	})
	require.NoError(t, err)

	require.Equal(t, repository.OrderStatusPendingPayment, resp.Order.Status)
	require.Equal(t, repository.OrderPaymentStatusPending, resp.Order.PaymentStatus)
	require.Equal(t, repository.PaymentMethodManual, resp.Order.PaymentMethod)
	require.Empty(t, resp.Order.PaymentIntentID)
	require.Nil(t, resp.Transaction)
	require.Empty(t, resp.Order.Payments)

	storedOrder, _, err := svcCtx.Repositories.Order.Get(ctx, resp.Order.ID)
	require.NoError(t, err)
	require.Equal(t, repository.OrderStatusPendingPayment, storedOrder.Status)
	require.Equal(t, repository.OrderPaymentStatusPending, storedOrder.PaymentStatus)
	require.Equal(t, repository.PaymentMethodManual, storedOrder.PaymentMethod)
}

func TestCreateOrderIdempotent(t *testing.T) {
	svcCtx, cleanup := setupCreateLogicTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC()

	user := repository.User{
		Email:       "buyer3@test.dev",
		DisplayName: "Buyer 3",
		Roles:       []string{"user"},
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	require.NoError(t, svcCtx.DB.Create(&user).Error)

	plan := repository.Plan{
		Name:              "Standard",
		Slug:              "standard",
		Description:       "Standard plan",
		PriceCents:        2000,
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

	balanceRepo := svcCtx.Repositories.Balance
	_, _, err := balanceRepo.ApplyTransaction(ctx, user.ID, repository.BalanceTransaction{
		Type:        "recharge",
		AmountCents: 10000,
		Currency:    "CNY",
		Reference:   "seed",
		Description: "seed balance",
	})
	require.NoError(t, err)

	claims := security.UserClaims{ID: user.ID, Email: user.Email, Roles: []string{"user"}}
	reqCtx := security.WithUser(ctx, claims)

	logic := NewCreateLogic(reqCtx, svcCtx)
	idemKey := "order-123"

	first, err := logic.Create(&types.UserCreateOrderRequest{
		PlanID:         plan.ID,
		Quantity:       2,
		IdempotencyKey: idemKey,
	})
	require.NoError(t, err)
	require.NotNil(t, first.Transaction)
	require.Equal(t, int64(6000), first.Balance.BalanceCents)

	second, err := logic.Create(&types.UserCreateOrderRequest{
		PlanID:         plan.ID,
		Quantity:       5, // should be ignored due to idempotency
		IdempotencyKey: idemKey,
	})
	require.NoError(t, err)
	require.Nil(t, second.Transaction)
	require.Equal(t, first.Order.ID, second.Order.ID)
	require.Equal(t, first.Order.Number, second.Order.Number)
	require.Equal(t, int64(6000), second.Balance.BalanceCents)
	require.Equal(t, 2, second.Order.Items[0].Quantity)

	var orderCount int64
	require.NoError(t, svcCtx.DB.Model(&repository.Order{}).Count(&orderCount).Error)
	require.EqualValues(t, 1, orderCount)

	txList, _, err := balanceRepo.ListTransactions(ctx, user.ID, repository.ListBalanceTransactionsOptions{Type: "purchase"})
	require.NoError(t, err)
	require.Len(t, txList, 1)
}
