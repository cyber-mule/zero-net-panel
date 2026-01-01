package order

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/logic/orderutil"
	"github.com/zero-net-panel/zero-net-panel/internal/logic/paymentutil"
	"github.com/zero-net-panel/zero-net-panel/internal/logic/subscriptionutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
	"github.com/zero-net-panel/zero-net-panel/pkg/metrics"
)

// CreateLogic handles user order creation.
type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic constructs CreateLogic.
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create issues an order for the given plan and settles payment according to the selected method.
func (l *CreateLogic) Create(req *types.UserCreateOrderRequest) (resp *types.UserOrderResponse, err error) {
	start := time.Now()
	idempotencyKey := strings.TrimSpace(req.IdempotencyKey)
	var idemPtr *string
	if idempotencyKey != "" {
		idemPtr = &idempotencyKey
	}
	method := strings.TrimSpace(strings.ToLower(req.PaymentMethod))
	if method == "" {
		method = repository.PaymentMethodBalance
	}
	if method == "offline" {
		method = repository.PaymentMethodManual
	}
	paymentMethod := method
	defer func() {
		result := "success"
		if err != nil {
			result = "error"
		}
		metrics.ObserveOrderCreate(paymentMethod, result, time.Since(start))
	}()

	if method != repository.PaymentMethodBalance && method != repository.PaymentMethodExternal && method != repository.PaymentMethodManual {
		return nil, repository.ErrInvalidArgument
	}

	user, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}

	if idempotencyKey != "" {
		if existing, items, payments, err := l.svcCtx.Repositories.Order.GetByIdempotencyKey(l.ctx, user.ID, idempotencyKey); err == nil {
			balance, balErr := l.svcCtx.Repositories.Balance.GetBalance(l.ctx, user.ID)
			if balErr != nil {
				return nil, balErr
			}
			detail := orderutil.ToOrderDetail(existing, items, nil, payments)
			resp := &types.UserOrderResponse{
				Order:   detail,
				Balance: orderutil.ToBalanceSnapshot(balance),
			}
			return resp, nil
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	if req.PlanID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	couponCode := strings.TrimSpace(req.CouponCode)

	plan, err := l.svcCtx.Repositories.Plan.Get(l.ctx, req.PlanID)
	if err != nil {
		return nil, err
	}

	if !plan.Visible || !strings.EqualFold(plan.Status, "active") {
		return nil, repository.ErrInvalidArgument
	}

	var billingOption repository.PlanBillingOption
	hasBillingOption := false
	billingOptionID := req.BillingOptionID
	if billingOptionID > 0 {
		option, err := l.svcCtx.Repositories.PlanBillingOption.Get(l.ctx, billingOptionID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, repository.ErrInvalidArgument
			}
			return nil, err
		}
		if option.PlanID != plan.ID {
			return nil, repository.ErrInvalidArgument
		}
		if !option.Visible || !strings.EqualFold(option.Status, "active") {
			return nil, repository.ErrInvalidArgument
		}
		billingOption = option
		hasBillingOption = true
	}

	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}
	if quantity > 10 {
		quantity = 10
	}

	channel := strings.TrimSpace(strings.ToLower(req.PaymentChannel))
	returnURL := strings.TrimSpace(req.PaymentReturnURL)
	var paymentChannel repository.PaymentChannel

	unitPriceCents := plan.PriceCents
	durationValue := plan.DurationDays
	durationUnit := repository.DurationUnitDay
	billingOptionName := ""
	if hasBillingOption {
		unitPriceCents = billingOption.PriceCents
		durationValue = billingOption.DurationValue
		durationUnit = strings.TrimSpace(strings.ToLower(billingOption.DurationUnit))
		if durationUnit == "" {
			durationUnit = repository.DurationUnitDay
		}
		switch durationUnit {
		case repository.DurationUnitHour, repository.DurationUnitDay, repository.DurationUnitMonth, repository.DurationUnitYear:
		default:
			return nil, repository.ErrInvalidArgument
		}
		billingOptionName = strings.TrimSpace(billingOption.Name)
	}

	totalCents := unitPriceCents * int64(quantity)
	if method == repository.PaymentMethodExternal && totalCents > 0 {
		if channel == "" {
			return nil, repository.ErrInvalidArgument
		}
		paymentChannel, err = l.svcCtx.Repositories.PaymentChannel.GetByCode(l.ctx, channel)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, repository.ErrInvalidArgument
			}
			return nil, err
		}
		if !paymentChannel.Enabled {
			return nil, repository.ErrInvalidArgument
		}
		channel = paymentChannel.Code
	}

	orderNumber := repository.GenerateOrderNumber()

	var createdOrder repository.Order
	var createdItems []repository.OrderItem
	var createdPayments []repository.OrderPayment
	var balance repository.UserBalance
	var balanceTx repository.BalanceTransaction
	var appliedCoupon *repository.Coupon
	var discountCents int64

	err = l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		orderRepo, err := repository.NewOrderRepository(tx)
		if err != nil {
			return err
		}
		balanceRepo, err := repository.NewBalanceRepository(tx)
		if err != nil {
			return err
		}
		couponRepo, err := repository.NewCouponRepository(tx)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		existingBalance, err := balanceRepo.GetBalance(l.ctx, user.ID)
		if err != nil {
			return err
		}
		balance = existingBalance

		currency := strings.TrimSpace(plan.Currency)
		if hasBillingOption {
			if optionCurrency := strings.TrimSpace(billingOption.Currency); optionCurrency != "" {
				currency = optionCurrency
			}
		}
		if currency == "" {
			currency = strings.TrimSpace(balance.Currency)
			if currency == "" {
				currency = "CNY"
			}
		}

		bindingIDs, err := l.svcCtx.Repositories.PlanProtocolBinding.ListBindingIDs(l.ctx, plan.ID)
		if err != nil {
			return err
		}
		bindingIDs = uniqueUint64s(bindingIDs)

		snapshot := map[string]any{
			"id":                  plan.ID,
			"name":                plan.Name,
			"slug":                plan.Slug,
			"description":         plan.Description,
			"price_cents":         unitPriceCents,
			"currency":            currency,
			"duration_unit":       durationUnit,
			"duration_value":      durationValue,
			"traffic_limit_bytes": plan.TrafficLimitBytes,
			"traffic_multipliers": cloneTrafficMultipliers(plan.TrafficMultipliers),
			"devices_limit":       plan.DevicesLimit,
			"features":            plan.Features,
			"tags":                plan.Tags,
		}
		if len(bindingIDs) > 0 {
			snapshot["binding_ids"] = bindingIDs
		}
		if durationUnit == repository.DurationUnitDay && durationValue > 0 {
			snapshot["duration_days"] = durationValue
		}
		if hasBillingOption {
			snapshot["billing_option_id"] = billingOption.ID
			if billingOptionName != "" {
				snapshot["billing_option_name"] = billingOptionName
			}
		}

		metadata := map[string]any{
			"quantity": quantity,
		}
		if channel != "" {
			metadata["payment_channel"] = channel
		}
		if returnURL != "" {
			metadata["payment_return_url"] = returnURL
		}

		baseTotalCents := unitPriceCents * int64(quantity)
		totalCents := baseTotalCents

		if couponCode != "" {
			coupon, err := couponRepo.GetByCodeForUpdate(l.ctx, couponCode)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return repository.ErrInvalidArgument
				}
				return err
			}
			if !strings.EqualFold(coupon.Status, repository.CouponStatusActive) {
				return repository.ErrInvalidArgument
			}
			if !coupon.StartsAt.IsZero() && now.Before(coupon.StartsAt) {
				return repository.ErrInvalidArgument
			}
			if !coupon.EndsAt.IsZero() && now.After(coupon.EndsAt) {
				return repository.ErrInvalidArgument
			}
			if coupon.MinOrderCents > 0 && baseTotalCents < coupon.MinOrderCents {
				return repository.ErrInvalidArgument
			}
			if coupon.MaxRedemptions > 0 {
				count, err := couponRepo.CountRedemptions(l.ctx, coupon.ID)
				if err != nil {
					return err
				}
				if count >= int64(coupon.MaxRedemptions) {
					return repository.ErrInvalidArgument
				}
			}
			if coupon.MaxRedemptionsPerUser > 0 {
				count, err := couponRepo.CountRedemptionsByUser(l.ctx, coupon.ID, user.ID)
				if err != nil {
					return err
				}
				if count >= int64(coupon.MaxRedemptionsPerUser) {
					return repository.ErrInvalidArgument
				}
			}

			amount, err := calculateDiscount(coupon, baseTotalCents, currency)
			if err != nil {
				return err
			}
			if amount > totalCents {
				amount = totalCents
			}
			if amount <= 0 {
				return repository.ErrInvalidArgument
			}

			appliedCoupon = &coupon
			discountCents = amount
			totalCents -= amount
			metadata["coupon_code"] = coupon.Code
			metadata["coupon_id"] = coupon.ID
			metadata["discount_cents"] = discountCents
		}

		orderModel := repository.Order{
			Number:         orderNumber,
			UserID:         user.ID,
			IdempotencyKey: idemPtr,
			PlanID:         &plan.ID,
			Status:         repository.OrderStatusPendingPayment,
			PaymentMethod:  method,
			PaymentStatus:  repository.OrderPaymentStatusPending,
			TotalCents:     totalCents,
			Currency:       currency,
			Metadata:       metadata,
			PlanSnapshot:   snapshot,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if totalCents == 0 || method == repository.PaymentMethodBalance {
			if totalCents > 0 {
				txRecord := repository.BalanceTransaction{
					Type:        "purchase",
					AmountCents: -totalCents,
					Currency:    currency,
					Reference:   fmt.Sprintf("order:%s", orderNumber),
					Description: fmt.Sprintf("购买套餐 %s", plan.Name),
					Metadata: map[string]any{
						"plan_id":      plan.ID,
						"quantity":     quantity,
						"order_number": orderNumber,
					},
				}
				if hasBillingOption {
					txRecord.Metadata["billing_option_id"] = billingOption.ID
				}
				createdTx, updatedBalance, err := balanceRepo.ApplyTransaction(l.ctx, user.ID, txRecord)
				if err != nil {
					return err
				}
				balanceTx = createdTx
				balance = updatedBalance
				paidAt := createdTx.CreatedAt.UTC()
				orderModel.Status = repository.OrderStatusPaid
				orderModel.PaymentStatus = repository.OrderPaymentStatusSucceeded
				orderModel.PaidAt = &paidAt
			} else {
				paidAt := now
				orderModel.Status = repository.OrderStatusPaid
				orderModel.PaymentStatus = repository.OrderPaymentStatusSucceeded
				orderModel.PaidAt = &paidAt
			}
		} else if method == repository.PaymentMethodExternal {
			intentID := fmt.Sprintf("%s-%s", channel, orderNumber)
			if channel == "" {
				intentID = orderNumber
			}
			orderModel.PaymentIntentID = intentID
		}

		item := repository.OrderItem{
			ItemType:       "plan",
			ItemID:         plan.ID,
			Name:           plan.Name,
			Quantity:       quantity,
			UnitPriceCents: unitPriceCents,
			Currency:       currency,
			SubtotalCents:  baseTotalCents,
			Metadata: map[string]any{
				"duration_unit":       durationUnit,
				"duration_value":      durationValue,
				"traffic_limit_bytes": plan.TrafficLimitBytes,
				"devices_limit":       plan.DevicesLimit,
			},
			CreatedAt: now,
		}
		if durationUnit == repository.DurationUnitDay && durationValue > 0 {
			item.Metadata["duration_days"] = durationValue
		}
		if hasBillingOption {
			item.Metadata["billing_option_id"] = billingOption.ID
			if billingOptionName != "" {
				item.Metadata["billing_option_name"] = billingOptionName
			}
		}

		itemsToCreate := []repository.OrderItem{item}
		if appliedCoupon != nil && discountCents > 0 {
			itemsToCreate = append(itemsToCreate, repository.OrderItem{
				ItemType:       "discount",
				ItemID:         appliedCoupon.ID,
				Name:           fmt.Sprintf("Coupon %s", appliedCoupon.Code),
				Quantity:       1,
				UnitPriceCents: -discountCents,
				Currency:       currency,
				SubtotalCents:  -discountCents,
				Metadata: map[string]any{
					"coupon_id":      appliedCoupon.ID,
					"coupon_code":    appliedCoupon.Code,
					"discount_type":  appliedCoupon.DiscountType,
					"discount_value": appliedCoupon.DiscountValue,
				},
				CreatedAt: now,
			})
		}

		created, items, err := orderRepo.Create(l.ctx, orderModel, itemsToCreate)
		if err != nil {
			return err
		}
		createdOrder = created
		createdItems = items

		if method == repository.PaymentMethodExternal && totalCents > 0 {
			paymentMetadata := map[string]any{}
			if channel != "" {
				paymentMetadata["channel"] = channel
			}
			if returnURL != "" {
				paymentMetadata["return_url"] = returnURL
			}
			paymentRecord := repository.OrderPayment{
				OrderID:     created.ID,
				Provider:    channel,
				Method:      method,
				IntentID:    created.PaymentIntentID,
				Status:      repository.OrderPaymentStatusPending,
				AmountCents: totalCents,
				Currency:    currency,
				Metadata:    paymentMetadata,
			}
			payment, err := orderRepo.CreatePayment(l.ctx, paymentRecord)
			if err != nil {
				return err
			}
			createdPayments = append(createdPayments, payment)
		}

		if strings.EqualFold(createdOrder.Status, repository.OrderStatusPaid) &&
			strings.EqualFold(createdOrder.PaymentStatus, repository.OrderPaymentStatusSucceeded) {
			txRepos, err := repository.NewRepositories(tx)
			if err != nil {
				return err
			}
			provisioned, err := subscriptionutil.EnsureOrderSubscription(l.ctx, txRepos, createdOrder, createdItems)
			if err != nil {
				return err
			}
			createdOrder = provisioned.Order
		}

		if appliedCoupon != nil && discountCents > 0 {
			status := repository.CouponRedemptionReserved
			if strings.EqualFold(createdOrder.Status, repository.OrderStatusPaid) {
				status = repository.CouponRedemptionApplied
			}
			_, err := couponRepo.CreateRedemption(l.ctx, repository.CouponRedemption{
				CouponID:    appliedCoupon.ID,
				UserID:      user.ID,
				OrderID:     createdOrder.ID,
				Status:      status,
				AmountCents: discountCents,
				Currency:    currency,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) && idempotencyKey != "" {
			existing, items, payments, fetchErr := l.svcCtx.Repositories.Order.GetByIdempotencyKey(l.ctx, user.ID, idempotencyKey)
			if fetchErr == nil {
				balance, balErr := l.svcCtx.Repositories.Balance.GetBalance(l.ctx, user.ID)
				if balErr != nil {
					return nil, balErr
				}
				detail := orderutil.ToOrderDetail(existing, items, nil, payments)
				return &types.UserOrderResponse{
					Order:   detail,
					Balance: orderutil.ToBalanceSnapshot(balance),
				}, nil
			}
		}
		return nil, err
	}

	if method == repository.PaymentMethodExternal && totalCents > 0 && len(createdPayments) > 0 {
		initResult, err := paymentutil.Initiate(l.ctx, paymentutil.InitiateParams{
			Channel:   paymentChannel,
			Order:     createdOrder,
			Payment:   createdPayments[0],
			Quantity:  quantity,
			PlanID:    plan.ID,
			PlanName:  plan.Name,
			ReturnURL: returnURL,
		})
		if err != nil {
			return nil, err
		}
		if len(initResult.Metadata) > 0 || initResult.Reference != "" {
			updateParams := repository.UpdateOrderPaymentParams{
				Status:        repository.OrderPaymentStatusPending,
				MetadataPatch: initResult.Metadata,
			}
			if initResult.Reference != "" {
				ref := initResult.Reference
				updateParams.Reference = &ref
			}
			updatedPayment, err := l.svcCtx.Repositories.Order.UpdatePaymentRecord(l.ctx, createdPayments[0].ID, updateParams)
			if err != nil {
				return nil, err
			}
			createdPayments[0] = updatedPayment
		}
	}

	detail := orderutil.ToOrderDetail(createdOrder, createdItems, nil, createdPayments)
	balanceView := orderutil.ToBalanceSnapshot(balance)

	var txView *types.BalanceTransactionSummary
	if balanceTx.ID != 0 {
		summary := orderutil.ToBalanceTransactionView(balanceTx)
		txView = &summary
	}

	resp = &types.UserOrderResponse{
		Order:       detail,
		Balance:     balanceView,
		Transaction: txView,
	}

	return resp, nil
}

func calculateDiscount(coupon repository.Coupon, baseTotalCents int64, currency string) (int64, error) {
	if baseTotalCents <= 0 {
		return 0, repository.ErrInvalidArgument
	}
	switch strings.ToLower(strings.TrimSpace(coupon.DiscountType)) {
	case repository.CouponTypePercent:
		if coupon.DiscountValue <= 0 || coupon.DiscountValue > 10000 {
			return 0, repository.ErrInvalidArgument
		}
		return baseTotalCents * coupon.DiscountValue / 10000, nil
	case repository.CouponTypeFixed:
		if coupon.DiscountValue <= 0 {
			return 0, repository.ErrInvalidArgument
		}
		if strings.TrimSpace(coupon.Currency) != "" && !strings.EqualFold(coupon.Currency, currency) {
			return 0, repository.ErrInvalidArgument
		}
		return coupon.DiscountValue, nil
	default:
		return 0, repository.ErrInvalidArgument
	}
}

func uniqueUint64s(input []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(input))
	result := make([]uint64, 0, len(input))
	for _, value := range input {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func cloneTrafficMultipliers(input map[string]float64) map[string]float64 {
	if input == nil {
		return map[string]float64{}
	}
	result := make(map[string]float64, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}
