package coupons

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles coupon creation.
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

// Create creates a new coupon.
func (l *CreateLogic) Create(req *types.AdminCreateCouponRequest) (*types.CouponSummary, error) {
	status := normalizeStatus(req.Status)
	if status == "" {
		status = repository.CouponStatusActive
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if err := validateDiscount(req.DiscountType, req.DiscountValue, currency); err != nil {
		return nil, err
	}

	var startsAt time.Time
	if req.StartsAt != nil {
		value, err := parseOptionalTime(req.StartsAt)
		if err != nil {
			return nil, err
		}
		if value != nil {
			startsAt = *value
		}
	}
	var endsAt time.Time
	if req.EndsAt != nil {
		value, err := parseOptionalTime(req.EndsAt)
		if err != nil {
			return nil, err
		}
		if value != nil {
			endsAt = *value
		}
	}
	if !startsAt.IsZero() && !endsAt.IsZero() && endsAt.Before(startsAt) {
		return nil, repository.ErrInvalidArgument
	}

	maxRedemptions := 0
	if req.MaxRedemptions != nil {
		maxRedemptions = *req.MaxRedemptions
	}
	maxPerUser := 0
	if req.MaxRedemptionsPerUser != nil {
		maxPerUser = *req.MaxRedemptionsPerUser
	}
	minOrder := int64(0)
	if req.MinOrderCents != nil {
		minOrder = *req.MinOrderCents
	}

	coupon := repository.Coupon{
		Code:                  req.Code,
		Name:                  req.Name,
		Description:           strings.TrimSpace(req.Description),
		Status:                status,
		DiscountType:          req.DiscountType,
		DiscountValue:         req.DiscountValue,
		Currency:              currency,
		MaxRedemptions:        maxRedemptions,
		MaxRedemptionsPerUser: maxPerUser,
		MinOrderCents:         minOrder,
		StartsAt:              startsAt,
		EndsAt:                endsAt,
	}

	created, err := l.svcCtx.Repositories.Coupon.Create(l.ctx, coupon)
	if err != nil {
		return nil, err
	}

	summary := toCouponSummary(created)
	return &summary, nil
}
