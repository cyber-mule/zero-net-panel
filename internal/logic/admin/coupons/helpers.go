package coupons

import (
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toCouponSummary(coupon repository.Coupon) types.CouponSummary {
	return types.CouponSummary{
		ID:                    coupon.ID,
		Code:                  coupon.Code,
		Name:                  coupon.Name,
		Description:           coupon.Description,
		Status:                coupon.Status,
		DiscountType:          coupon.DiscountType,
		DiscountValue:         coupon.DiscountValue,
		Currency:              coupon.Currency,
		MaxRedemptions:        coupon.MaxRedemptions,
		MaxRedemptionsPerUser: coupon.MaxRedemptionsPerUser,
		MinOrderCents:         coupon.MinOrderCents,
		StartsAt:              toUnixPtr(coupon.StartsAt),
		EndsAt:                toUnixPtr(coupon.EndsAt),
		CreatedAt:             toUnixOrZero(coupon.CreatedAt),
		UpdatedAt:             toUnixOrZero(coupon.UpdatedAt),
	}
}

func normalizeStatus(statusCode int) (int, error) {
	switch statusCode {
	case repository.CouponStatusActive, repository.CouponStatusDisabled:
		return statusCode, nil
	case 0:
		return 0, repository.ErrInvalidArgument
	default:
		return 0, repository.ErrInvalidArgument
	}
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}

func toUnixPtr(ts time.Time) *int64 {
	if ts.IsZero() {
		return nil
	}
	value := ts.Unix()
	return &value
}

func parseOptionalTime(value *int64) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	if *value <= 0 {
		return nil, repository.ErrInvalidArgument
	}
	ts := time.Unix(*value, 0).UTC()
	return &ts, nil
}

func validateDiscount(discountType string, discountValue int64, currency string) error {
	discountType = strings.ToLower(strings.TrimSpace(discountType))
	if discountType == "" || discountValue <= 0 {
		return repository.ErrInvalidArgument
	}
	switch discountType {
	case repository.CouponTypePercent:
		if discountValue > 10000 {
			return repository.ErrInvalidArgument
		}
	case repository.CouponTypeFixed:
		if strings.TrimSpace(currency) == "" {
			return repository.ErrInvalidArgument
		}
	default:
		return repository.ErrInvalidArgument
	}
	return nil
}

func normalizePage(page, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}
