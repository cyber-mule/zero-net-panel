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

// UpdateLogic handles coupon updates.
type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLogic constructs UpdateLogic.
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update updates a coupon.
func (l *UpdateLogic) Update(req *types.AdminUpdateCouponRequest) (*types.CouponSummary, error) {
	if req.CouponID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	existing, err := l.svcCtx.Repositories.Coupon.Get(l.ctx, req.CouponID)
	if err != nil {
		return nil, err
	}

	discountType := existing.DiscountType
	if req.DiscountType != nil {
		discountType = *req.DiscountType
	}
	discountValue := existing.DiscountValue
	if req.DiscountValue != nil {
		discountValue = *req.DiscountValue
	}
	currency := existing.Currency
	if req.Currency != nil {
		currency = strings.ToUpper(strings.TrimSpace(*req.Currency))
	}
	if req.DiscountType != nil || req.DiscountValue != nil || req.Currency != nil {
		if err := validateDiscount(discountType, discountValue, currency); err != nil {
			return nil, err
		}
	}

	var startsAt *time.Time
	if req.StartsAt != nil {
		value, err := parseOptionalTime(req.StartsAt)
		if err != nil {
			return nil, err
		}
		startsAt = value
	}
	var endsAt *time.Time
	if req.EndsAt != nil {
		value, err := parseOptionalTime(req.EndsAt)
		if err != nil {
			return nil, err
		}
		endsAt = value
	}

	nextStarts := existing.StartsAt
	if startsAt != nil {
		nextStarts = *startsAt
	}
	nextEnds := existing.EndsAt
	if endsAt != nil {
		nextEnds = *endsAt
	}
	if !nextStarts.IsZero() && !nextEnds.IsZero() && nextEnds.Before(nextStarts) {
		return nil, repository.ErrInvalidArgument
	}

	var statusPtr *int
	if req.Status != nil {
		normalized, err := normalizeStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		statusPtr = &normalized
	}

	input := repository.UpdateCouponInput{
		Name:                  req.Name,
		Description:           req.Description,
		Status:                statusPtr,
		DiscountType:          req.DiscountType,
		DiscountValue:         req.DiscountValue,
		Currency:              req.Currency,
		MaxRedemptions:        req.MaxRedemptions,
		MaxRedemptionsPerUser: req.MaxRedemptionsPerUser,
		MinOrderCents:         req.MinOrderCents,
		StartsAt:              startsAt,
		EndsAt:                endsAt,
	}

	updated, err := l.svcCtx.Repositories.Coupon.Update(l.ctx, req.CouponID, input)
	if err != nil {
		return nil, err
	}

	summary := toCouponSummary(updated)
	return &summary, nil
}
