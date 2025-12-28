package coupons

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles coupon listing.
type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListLogic constructs ListLogic.
func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// List returns coupon list.
func (l *ListLogic) List(req *types.AdminListCouponsRequest) (*types.AdminCouponListResponse, error) {
	opts := repository.ListCouponsOptions{
		Page:      req.Page,
		PerPage:   req.PerPage,
		Sort:      req.Sort,
		Direction: req.Direction,
		Query:     req.Query,
		Status:    req.Status,
	}

	coupons, total, err := l.svcCtx.Repositories.Coupon.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	list := make([]types.CouponSummary, 0, len(coupons))
	for _, coupon := range coupons {
		list = append(list, toCouponSummary(coupon))
	}

	page, perPage := normalizePage(req.Page, req.PerPage)
	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.AdminCouponListResponse{
		Coupons:    list,
		Pagination: pagination,
	}, nil
}
