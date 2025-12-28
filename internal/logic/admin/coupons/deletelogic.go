package coupons

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// DeleteLogic handles coupon deletion.
type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteLogic constructs DeleteLogic.
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Delete removes a coupon.
func (l *DeleteLogic) Delete(req *types.AdminDeleteCouponRequest) error {
	if req.CouponID == 0 {
		return repository.ErrInvalidArgument
	}
	return l.svcCtx.Repositories.Coupon.Delete(l.ctx, req.CouponID)
}
