package orders

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/logic/orderutil"
	"github.com/zero-net-panel/zero-net-panel/internal/logic/subscriptionutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// PaymentCallbackLogic processes external payment notifications and updates order state.
type PaymentCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPaymentCallbackLogic constructs a new callback handler instance.
func NewPaymentCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentCallbackLogic {
	return &PaymentCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Process applies payment status updates from external gateways and returns the refreshed order state.
func (l *PaymentCallbackLogic) Process(req *types.AdminPaymentCallbackRequest) (*types.AdminOrderResponse, error) {
	if req == nil || req.OrderID == 0 || req.PaymentID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	status := req.Status
	if status == 0 {
		return nil, repository.ErrInvalidArgument
	}
	if status != repository.OrderPaymentStatusSucceeded && status != repository.OrderPaymentStatusFailed {
		return nil, repository.ErrInvalidArgument
	}

	order, items, err := l.svcCtx.Repositories.Order.Get(l.ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	paymentsMap, err := l.svcCtx.Repositories.Order.ListPayments(l.ctx, []uint64{order.ID})
	if err != nil {
		return nil, err
	}

	var existingPayment repository.OrderPayment
	found := false
	for _, payment := range paymentsMap[order.ID] {
		if payment.ID == req.PaymentID {
			existingPayment = payment
			found = true
			break
		}
	}
	if !found {
		return nil, repository.ErrNotFound
	}
	if existingPayment.Status == status {
		if status == repository.OrderPaymentStatusSucceeded {
			if err := l.svcCtx.Repositories.Coupon.UpdateRedemptionStatusByOrder(l.ctx, order.ID, repository.CouponRedemptionApplied); err != nil {
				return nil, err
			}
		} else if status == repository.OrderPaymentStatusFailed {
			if err := l.svcCtx.Repositories.Coupon.UpdateRedemptionStatusByOrder(l.ctx, order.ID, repository.CouponRedemptionReleased); err != nil {
				return nil, err
			}
		}

		provisionedOrder := order
		if status == repository.OrderPaymentStatusSucceeded && order.Status == repository.OrderStatusPaid {
			if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
				result, err := subscriptionutil.EnsureOrderSubscription(l.ctx, txRepos, order, items)
				if err != nil {
					return err
				}
				provisionedOrder = result.Order
				return nil
			}); err != nil {
				return nil, err
			}
		}

		refundsMap, err := l.svcCtx.Repositories.Order.ListRefunds(l.ctx, []uint64{order.ID})
		if err != nil {
			return nil, err
		}
		detail := orderutil.ToOrderDetail(provisionedOrder, items, refundsMap[order.ID], paymentsMap[order.ID])
		u, err := l.svcCtx.Repositories.User.Get(l.ctx, provisionedOrder.UserID)
		if err != nil {
			return nil, err
		}
		resp := types.AdminOrderResponse{
			Order: types.AdminOrderDetail{
				OrderDetail: detail,
				User: types.OrderUserSummary{
					ID:          u.ID,
					Email:       u.Email,
					DisplayName: u.DisplayName,
				},
			},
		}
		return &resp, nil
	}
	if existingPayment.Status == repository.OrderPaymentStatusSucceeded && status == repository.OrderPaymentStatusFailed {
		return nil, repository.ErrInvalidState
	}

	var updatedOrder repository.Order
	var updatedPayment repository.OrderPayment

	err = l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		repo, err := repository.NewOrderRepository(tx)
		if err != nil {
			return err
		}
		couponRepo, err := repository.NewCouponRepository(tx)
		if err != nil {
			return err
		}

		paymentParams := repository.UpdateOrderPaymentParams{Status: status}
		if ref := strings.TrimSpace(req.Reference); ref != "" {
			paymentParams.Reference = &ref
		}
		if code := strings.TrimSpace(req.FailureCode); code != "" {
			paymentParams.FailureCode = &code
		}
		if message := strings.TrimSpace(req.FailureMessage); message != "" {
			paymentParams.FailureMessage = &message
		}
		if req.PaidAt != nil && *req.PaidAt > 0 {
			processedAt := time.Unix(*req.PaidAt, 0).UTC()
			paymentParams.ProcessedAt = &processedAt
		}

		payment, err := repo.UpdatePaymentRecord(l.ctx, req.PaymentID, paymentParams)
		if err != nil {
			return err
		}
		updatedPayment = payment

		stateParams := repository.UpdateOrderPaymentStateParams{
			PaymentStatus: status,
		}

		if status == repository.OrderPaymentStatusSucceeded {
			orderStatus := repository.OrderStatusPaid
			stateParams.OrderStatus = &orderStatus
			paidAt := time.Now().UTC()
			if req.PaidAt != nil && *req.PaidAt > 0 {
				paidAt = time.Unix(*req.PaidAt, 0).UTC()
			}
			stateParams.PaidAt = &paidAt
			if ref := strings.TrimSpace(req.Reference); ref != "" {
				stateParams.PaymentReference = &ref
			}
		} else {
			orderStatus := repository.OrderStatusPaymentFailed
			stateParams.OrderStatus = &orderStatus
			if code := strings.TrimSpace(req.FailureCode); code != "" {
				stateParams.FailureCode = &code
			}
			if message := strings.TrimSpace(req.FailureMessage); message != "" {
				stateParams.FailureMessage = &message
			}
		}

		updated, err := repo.UpdatePaymentState(l.ctx, req.OrderID, stateParams)
		if err != nil {
			return err
		}
		updatedOrder = updated

		if status == repository.OrderPaymentStatusSucceeded {
			txRepos, err := repository.NewRepositories(tx)
			if err != nil {
				return err
			}
			provisioned, err := subscriptionutil.EnsureOrderSubscription(l.ctx, txRepos, updatedOrder, items)
			if err != nil {
				return err
			}
			updatedOrder = provisioned.Order
			if err := couponRepo.UpdateRedemptionStatusByOrder(l.ctx, updatedOrder.ID, repository.CouponRedemptionApplied); err != nil {
				return err
			}
		} else if status == repository.OrderPaymentStatusFailed {
			if err := couponRepo.UpdateRedemptionStatusByOrder(l.ctx, updatedOrder.ID, repository.CouponRedemptionReleased); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	payments := paymentsMap[order.ID]
	replaced := false
	for i := range payments {
		if payments[i].ID == updatedPayment.ID {
			payments[i] = updatedPayment
			replaced = true
			break
		}
	}
	if !replaced {
		payments = append(payments, updatedPayment)
	}

	refundsMap, err := l.svcCtx.Repositories.Order.ListRefunds(l.ctx, []uint64{order.ID})
	if err != nil {
		return nil, err
	}

	detail := orderutil.ToOrderDetail(updatedOrder, items, refundsMap[order.ID], payments)

	u, err := l.svcCtx.Repositories.User.Get(l.ctx, updatedOrder.UserID)
	if err != nil {
		return nil, err
	}

	resp := types.AdminOrderResponse{
		Order: types.AdminOrderDetail{
			OrderDetail: detail,
			User: types.OrderUserSummary{
				ID:          u.ID,
				Email:       u.Email,
				DisplayName: u.DisplayName,
			},
		},
	}

	l.Infof("audit: payment callback order=%d payment=%d status=%d reference=%s", updatedOrder.ID, updatedPayment.ID, status, strings.TrimSpace(req.Reference))

	return &resp, nil
}
