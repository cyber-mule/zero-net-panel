package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/bootstrap/migrations"
	adminorders "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/orders"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/testutil"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func setupTestServiceContext(t *testing.T) (*svc.ServiceContext, context.Context) {
	t.Helper()

	testutil.RequireSQLite(t)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	ctx := context.Background()
	if _, err := migrations.Apply(ctx, db, 0, false); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	repos, err := repository.NewRepositories(db)
	if err != nil {
		t.Fatalf("create repositories: %v", err)
	}

	svcCtx := &svc.ServiceContext{
		DB:           db,
		Repositories: repos,
	}

	return svcCtx, ctx
}

func seedDefaultSubscriptionTemplate(t *testing.T, db *gorm.DB) repository.SubscriptionTemplate {
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
		t.Fatalf("find template: %v", err)
	}
	if err := db.Create(&tpl).Error; err != nil {
		t.Fatalf("create template: %v", err)
	}
	return tpl
}

func TestOrderLifecycle(t *testing.T) {
	svcCtx, ctx := setupTestServiceContext(t)

	now := time.Now().UTC()

	customer := repository.User{
		Email:       "user@example.com",
		DisplayName: "Test User",
		Roles:       []string{"user"},
		Status:      status.UserStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := svcCtx.DB.Create(&customer).Error; err != nil {
		t.Fatalf("create customer: %v", err)
	}

	admin := repository.User{
		Email:       "admin@example.com",
		DisplayName: "Admin",
		Roles:       []string{"admin"},
		Status:      status.UserStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := svcCtx.DB.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	seedDefaultSubscriptionTemplate(t, svcCtx.DB)

	plan := repository.Plan{
		Name:              "Plan B",
		Slug:              "plan-b",
		Description:       "Plan B",
		PriceCents:        2000,
		Currency:          "CNY",
		DurationDays:      30,
		TrafficLimitBytes: 1024,
		DevicesLimit:      2,
		Status:            status.PlanStatusActive,
		Visible:           true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := svcCtx.DB.Create(&plan).Error; err != nil {
		t.Fatalf("create plan: %v", err)
	}

	orderRepo := svcCtx.Repositories.Order

	// Prepare order for user cancellation.
	cancelOrder, _, err := orderRepo.Create(ctx, repository.Order{
		UserID:        customer.ID,
		Status:        repository.OrderStatusPending,
		PaymentMethod: repository.PaymentMethodBalance,
		TotalCents:    1200,
		Currency:      "CNY",
		Metadata:      map[string]any{},
	}, []repository.OrderItem{{
		ItemType:       "plan",
		ItemID:         1,
		Name:           "Plan A",
		Quantity:       1,
		UnitPriceCents: 1200,
		Currency:       "CNY",
		SubtotalCents:  1200,
	}})
	if err != nil {
		t.Fatalf("create cancel order: %v", err)
	}

	// User cancels pending order.
	userCtx := security.WithUser(context.Background(), security.UserClaims{ID: customer.ID, Roles: []string{"user"}})
	cancelLogic := NewCancelLogic(userCtx, svcCtx)
	cancelResp, err := cancelLogic.Cancel(&types.UserCancelOrderRequest{OrderID: cancelOrder.ID, Reason: "no longer needed"})
	if err != nil {
		t.Fatalf("cancel order: %v", err)
	}
	if cancelResp.Order.Status != repository.OrderStatusCancelled {
		t.Fatalf("expected order status cancelled, got %d", cancelResp.Order.Status)
	}
	if cancelResp.Order.Metadata["cancel_reason"] != "no longer needed" {
		t.Fatalf("expected cancel reason metadata, got %v", cancelResp.Order.Metadata["cancel_reason"])
	}

	// Prepare second order for manual payment and refunds.
	payOrder, _, err := orderRepo.Create(ctx, repository.Order{
		UserID:        customer.ID,
		PlanID:        &plan.ID,
		Status:        repository.OrderStatusPending,
		PaymentMethod: repository.PaymentMethodBalance,
		TotalCents:    plan.PriceCents,
		Currency:      plan.Currency,
		PlanSnapshot: map[string]any{
			"name":                plan.Name,
			"duration_days":       plan.DurationDays,
			"traffic_limit_bytes": plan.TrafficLimitBytes,
			"devices_limit":       plan.DevicesLimit,
		},
		Metadata: map[string]any{},
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
	if err != nil {
		t.Fatalf("create pay order: %v", err)
	}

	balanceRepo := svcCtx.Repositories.Balance
	if _, _, err := balanceRepo.ApplyTransaction(ctx, customer.ID, repository.BalanceTransaction{
		Type:        "recharge",
		AmountCents: 5000,
		Currency:    "CNY",
		Reference:   "test-topup",
		Description: "seed balance",
	}); err != nil {
		t.Fatalf("seed balance: %v", err)
	}

	adminCtx := security.WithUser(context.Background(), security.UserClaims{ID: admin.ID, Roles: []string{"admin"}})

	markLogic := adminorders.NewMarkPaidLogic(adminCtx, svcCtx)
	markResp, err := markLogic.MarkPaid(&types.AdminMarkOrderPaidRequest{
		OrderID:       payOrder.ID,
		PaymentMethod: "manual",
		Note:          "manual charge",
		Reference:     "manual-ref",
		ChargeBalance: true,
	})
	if err != nil {
		t.Fatalf("mark order paid: %v", err)
	}
	if markResp.Order.Status != repository.OrderStatusPaid {
		t.Fatalf("expected paid status, got %d", markResp.Order.Status)
	}
	if markResp.Order.PaymentMethod != repository.PaymentMethodBalance {
		t.Fatalf("expected payment method balance, got %s", markResp.Order.PaymentMethod)
	}

	balanceAfterPay, err := balanceRepo.GetBalance(ctx, customer.ID)
	if err != nil {
		t.Fatalf("get balance after pay: %v", err)
	}
	if balanceAfterPay.BalanceCents != 3000 {
		t.Fatalf("expected balance 3000 after charge, got %d", balanceAfterPay.BalanceCents)
	}

	subs, _, err := svcCtx.Repositories.Subscription.ListByUser(ctx, customer.ID, repository.ListSubscriptionsOptions{})
	if err != nil {
		t.Fatalf("list subscriptions: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(subs))
	}
	if subs[0].PlanName != plan.Name {
		t.Fatalf("expected plan name %s, got %s", plan.Name, subs[0].PlanName)
	}

	refundLogic := adminorders.NewRefundLogic(adminCtx, svcCtx)
	half := payOrder.TotalCents / 2
	refundResp1, err := refundLogic.Refund(&types.AdminRefundOrderRequest{
		OrderID:       payOrder.ID,
		AmountCents:   half,
		Reason:        "partial",
		CreditBalance: true,
	})
	if err != nil {
		t.Fatalf("partial refund: %v", err)
	}
	if refundResp1.Order.Status != repository.OrderStatusPartiallyRefunded {
		t.Fatalf("expected partially_refunded status after partial refund, got %d", refundResp1.Order.Status)
	}
	if refundResp1.Order.RefundedCents != half {
		t.Fatalf("expected refunded cents %d, got %d", half, refundResp1.Order.RefundedCents)
	}

	balanceAfterPartial, err := balanceRepo.GetBalance(ctx, customer.ID)
	if err != nil {
		t.Fatalf("get balance after partial refund: %v", err)
	}
	if balanceAfterPartial.BalanceCents != 3000+half {
		t.Fatalf("expected balance %d after partial refund, got %d", 3000+half, balanceAfterPartial.BalanceCents)
	}

	finalAmount := payOrder.TotalCents - half
	refundResp2, err := refundLogic.Refund(&types.AdminRefundOrderRequest{
		OrderID:       payOrder.ID,
		AmountCents:   finalAmount,
		Reason:        "final",
		CreditBalance: true,
	})
	if err != nil {
		t.Fatalf("final refund: %v", err)
	}
	if refundResp2.Order.Status != repository.OrderStatusRefunded {
		t.Fatalf("expected refunded status after full refund, got %d", refundResp2.Order.Status)
	}
	if refundResp2.Order.RefundedCents != payOrder.TotalCents {
		t.Fatalf("expected refunded cents %d, got %d", payOrder.TotalCents, refundResp2.Order.RefundedCents)
	}
	if len(refundResp2.Order.Refunds) != 2 {
		t.Fatalf("expected 2 refund records, got %d", len(refundResp2.Order.Refunds))
	}

	balanceAfterFull, err := balanceRepo.GetBalance(ctx, customer.ID)
	if err != nil {
		t.Fatalf("get balance after full refund: %v", err)
	}
	if balanceAfterFull.BalanceCents != 5000 {
		t.Fatalf("expected balance restored to 5000, got %d", balanceAfterFull.BalanceCents)
	}

	refundsMap, err := orderRepo.ListRefunds(ctx, []uint64{payOrder.ID})
	if err != nil {
		t.Fatalf("list refunds: %v", err)
	}
	if len(refundsMap[payOrder.ID]) != 2 {
		t.Fatalf("expected repository to return 2 refunds, got %d", len(refundsMap[payOrder.ID]))
	}

	// Admin cancel another pending order.
	cancelAdminOrder, _, err := orderRepo.Create(ctx, repository.Order{
		UserID:        customer.ID,
		Status:        repository.OrderStatusPending,
		PaymentMethod: repository.PaymentMethodBalance,
		TotalCents:    1500,
		Currency:      "CNY",
		Metadata:      map[string]any{},
	}, []repository.OrderItem{{
		ItemType:       "plan",
		ItemID:         3,
		Name:           "Plan C",
		Quantity:       1,
		UnitPriceCents: 1500,
		Currency:       "CNY",
		SubtotalCents:  1500,
	}})
	if err != nil {
		t.Fatalf("create admin cancel order: %v", err)
	}

	adminCancelLogic := adminorders.NewCancelLogic(adminCtx, svcCtx)
	cancelRespAdmin, err := adminCancelLogic.Cancel(&types.AdminCancelOrderRequest{OrderID: cancelAdminOrder.ID, Reason: "fraud"})
	if err != nil {
		t.Fatalf("admin cancel: %v", err)
	}
	if cancelRespAdmin.Order.Status != repository.OrderStatusCancelled {
		t.Fatalf("expected admin cancelled status, got %d", cancelRespAdmin.Order.Status)
	}
}
