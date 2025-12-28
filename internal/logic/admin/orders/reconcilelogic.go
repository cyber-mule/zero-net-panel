package orders

import (
	"context"
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
)

// ReconcileLogic handles admin payment reconciliation.
type ReconcileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewReconcileLogic constructs ReconcileLogic.
func NewReconcileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReconcileLogic {
	return &ReconcileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Reconcile queries gateway payment status and updates the order state.
func (l *ReconcileLogic) Reconcile(req *types.AdminReconcilePaymentRequest) (*types.AdminOrderResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}
	if req == nil || req.PaymentID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	payment, err := l.svcCtx.Repositories.Order.GetPayment(l.ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}
	if req.OrderID != 0 && payment.OrderID != req.OrderID {
		return nil, repository.ErrInvalidArgument
	}

	order, items, err := l.svcCtx.Repositories.Order.Get(l.ctx, payment.OrderID)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(order.PaymentMethod, repository.PaymentMethodExternal) {
		return nil, repository.ErrInvalidArgument
	}

	channel, err := l.svcCtx.Repositories.PaymentChannel.GetByCode(l.ctx, payment.Provider)
	if err != nil {
		return nil, err
	}

	result, err := paymentutil.Reconcile(l.ctx, paymentutil.ReconcileParams{
		Channel: channel,
		Order:   order,
		Payment: payment,
	})
	if err != nil {
		return nil, err
	}

	if strings.EqualFold(payment.Status, repository.OrderPaymentStatusSucceeded) && result.Status == repository.OrderPaymentStatusFailed {
		return nil, repository.ErrInvalidState
	}

	metadataPatch := result.Metadata
	if metadataPatch == nil {
		metadataPatch = map[string]any{}
	}
	metadataPatch["reconciled_at"] = time.Now().UTC().Unix()

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

		paymentStatus := payment.Status
		if result.Status != repository.OrderPaymentStatusPending {
			paymentStatus = result.Status
		}
		if strings.EqualFold(payment.Status, repository.OrderPaymentStatusSucceeded) && paymentStatus == repository.OrderPaymentStatusPending {
			paymentStatus = payment.Status
		}

		paymentParams := repository.UpdateOrderPaymentParams{
			Status:        paymentStatus,
			MetadataPatch: metadataPatch,
		}
		if ref := strings.TrimSpace(result.Reference); ref != "" {
			paymentParams.Reference = &ref
		}
		if code := strings.TrimSpace(result.FailureCode); code != "" {
			paymentParams.FailureCode = &code
		}
		if message := strings.TrimSpace(result.FailureMessage); message != "" {
			paymentParams.FailureMessage = &message
		}

		updated, err := repo.UpdatePaymentRecord(l.ctx, payment.ID, paymentParams)
		if err != nil {
			return err
		}
		updatedPayment = updated
		updatedOrder = order

		if paymentStatus == repository.OrderPaymentStatusPending {
			return nil
		}

		stateParams := repository.UpdateOrderPaymentStateParams{
			PaymentStatus: paymentStatus,
		}

		if paymentStatus == repository.OrderPaymentStatusSucceeded {
			orderStatus := repository.OrderStatusPaid
			stateParams.OrderStatus = &orderStatus
			paidAt := time.Now().UTC()
			stateParams.PaidAt = &paidAt
			if ref := strings.TrimSpace(result.Reference); ref != "" {
				stateParams.PaymentReference = &ref
			}
		} else {
			orderStatus := repository.OrderStatusPaymentFailed
			stateParams.OrderStatus = &orderStatus
			if code := strings.TrimSpace(result.FailureCode); code != "" {
				stateParams.FailureCode = &code
			}
			if message := strings.TrimSpace(result.FailureMessage); message != "" {
				stateParams.FailureMessage = &message
			}
		}

		updatedOrder, err = repo.UpdatePaymentState(l.ctx, order.ID, stateParams)
		if err != nil {
			return err
		}

		if paymentStatus == repository.OrderPaymentStatusSucceeded {
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
		} else if paymentStatus == repository.OrderPaymentStatusFailed {
			if err := couponRepo.UpdateRedemptionStatusByOrder(l.ctx, updatedOrder.ID, repository.CouponRedemptionReleased); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	paymentsMap, err := l.svcCtx.Repositories.Order.ListPayments(l.ctx, []uint64{order.ID})
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

	l.Infof("audit: payment reconcile order=%d payment=%d status=%s actor=%s", updatedOrder.ID, updatedPayment.ID, updatedPayment.Status, actor.Email)

	return &resp, nil
}
