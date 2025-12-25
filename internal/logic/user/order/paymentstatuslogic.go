package order

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// PaymentStatusLogic handles retrieving payment status for a user order.
type PaymentStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPaymentStatusLogic constructs PaymentStatusLogic.
func NewPaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentStatusLogic {
	return &PaymentStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Get returns payment status for an order owned by the user.
func (l *PaymentStatusLogic) Get(req *types.UserOrderPaymentStatusRequest) (*types.UserOrderPaymentStatusResponse, error) {
	user, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}

	order, _, err := l.svcCtx.Repositories.Order.Get(l.ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != user.ID {
		return nil, repository.ErrForbidden
	}

	resp := &types.UserOrderPaymentStatusResponse{
		OrderID:       order.ID,
		Status:        order.Status,
		PaymentStatus: order.PaymentStatus,
		PaymentMethod: order.PaymentMethod,
		RefundedCents: order.RefundedCents,
		UpdatedAt:     order.UpdatedAt.UTC().Unix(),
	}

	if order.PaymentIntentID != "" {
		intent := order.PaymentIntentID
		resp.PaymentIntentID = &intent
	}
	if order.PaymentReference != "" {
		ref := order.PaymentReference
		resp.PaymentReference = &ref
	}
	if order.PaymentFailureCode != "" {
		code := order.PaymentFailureCode
		resp.PaymentFailureCode = &code
	}
	if order.PaymentFailureReason != "" {
		msg := order.PaymentFailureReason
		resp.PaymentFailureMessage = &msg
	}

	resp.PaidAt = toUnixPtr(order.PaidAt)
	resp.CancelledAt = toUnixPtr(order.CancelledAt)
	resp.RefundedAt = toUnixPtr(order.RefundedAt)

	return resp, nil
}

func toUnixPtr(ts *time.Time) *int64 {
	if ts == nil {
		return nil
	}
	value := ts.UTC().Unix()
	return &value
}
