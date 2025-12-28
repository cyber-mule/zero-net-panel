package types

// CouponSummary 优惠券摘要。
type CouponSummary struct {
	ID                    uint64 `json:"id"`
	Code                  string `json:"code"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Status                string `json:"status"`
	DiscountType          string `json:"discount_type"`
	DiscountValue         int64  `json:"discount_value"`
	Currency              string `json:"currency"`
	MaxRedemptions        int    `json:"max_redemptions"`
	MaxRedemptionsPerUser int    `json:"max_redemptions_per_user"`
	MinOrderCents         int64  `json:"min_order_cents"`
	StartsAt              *int64 `json:"starts_at,omitempty"`
	EndsAt                *int64 `json:"ends_at,omitempty"`
	CreatedAt             int64  `json:"created_at"`
	UpdatedAt             int64  `json:"updated_at"`
}

// AdminListCouponsRequest 管理端优惠券列表请求。
type AdminListCouponsRequest struct {
	Page      int    `form:"page,optional" json:"page,optional"`
	PerPage   int    `form:"per_page,optional" json:"per_page,optional"`
	Query     string `form:"q,optional" json:"q,optional"`
	Status    string `form:"status,optional" json:"status,optional"`
	Sort      string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
}

// AdminCouponListResponse 管理端优惠券列表响应。
type AdminCouponListResponse struct {
	Coupons    []CouponSummary `json:"coupons"`
	Pagination PaginationMeta  `json:"pagination"`
}

// AdminCreateCouponRequest 管理端创建优惠券请求。
type AdminCreateCouponRequest struct {
	Code                  string `json:"code"`
	Name                  string `json:"name"`
	Description           string `json:"description,omitempty,optional"`
	Status                string `json:"status,omitempty,optional"`
	DiscountType          string `json:"discount_type"`
	DiscountValue         int64  `json:"discount_value"`
	Currency              string `json:"currency,omitempty,optional"`
	MaxRedemptions        *int   `json:"max_redemptions,omitempty,optional"`
	MaxRedemptionsPerUser *int   `json:"max_redemptions_per_user,omitempty,optional"`
	MinOrderCents         *int64 `json:"min_order_cents,omitempty,optional"`
	StartsAt              *int64 `json:"starts_at,omitempty,optional"`
	EndsAt                *int64 `json:"ends_at,omitempty,optional"`
}

// AdminUpdateCouponRequest 管理端更新优惠券请求。
type AdminUpdateCouponRequest struct {
	CouponID              uint64  `path:"id"`
	Name                  *string `json:"name,omitempty,optional"`
	Description           *string `json:"description,omitempty,optional"`
	Status                *string `json:"status,omitempty,optional"`
	DiscountType          *string `json:"discount_type,omitempty,optional"`
	DiscountValue         *int64  `json:"discount_value,omitempty,optional"`
	Currency              *string `json:"currency,omitempty,optional"`
	MaxRedemptions        *int    `json:"max_redemptions,omitempty,optional"`
	MaxRedemptionsPerUser *int    `json:"max_redemptions_per_user,omitempty,optional"`
	MinOrderCents         *int64  `json:"min_order_cents,omitempty,optional"`
	StartsAt              *int64  `json:"starts_at,omitempty,optional"`
	EndsAt                *int64  `json:"ends_at,omitempty,optional"`
}

// AdminDeleteCouponRequest 管理端删除优惠券请求。
type AdminDeleteCouponRequest struct {
	CouponID uint64 `path:"id"`
}
