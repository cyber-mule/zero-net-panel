### 1. "List coupons"

1. route definition

- Url: /api/v1/admin/coupons
- Method: GET
- Request: `AdminListCouponsRequest`
- Response: `AdminCouponListResponse`

2. request definition



```golang
type AdminListCouponsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
}
```


3. response definition



```golang
type AdminCouponListResponse struct {
	Coupons []CouponSummary 
	Pagination PaginationMeta 
}

type PaginationMeta struct {
	Page int 
	Per_page int 
	Total_count int64 
	Has_next bool 
	Has_prev bool 
}
```

### 2. "Create coupon"

1. route definition

- Url: /api/v1/admin/coupons
- Method: POST
- Request: `AdminCreateCouponRequest`
- Response: `CouponSummary`

2. request definition



```golang
type AdminCreateCouponRequest struct {
	Code string 
	Name string 
	Description string `form:"description,optional" json:"description,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Discount_type string 
	Discount_value int64 
	Currency string `form:"currency,optional" json:"currency,optional"`
	Max_redemptions int `form:"max_redemptions,optional" json:"max_redemptions,optional"`
	Max_redemptions_per_user int `form:"max_redemptions_per_user,optional" json:"max_redemptions_per_user,optional"`
	Min_order_cents int64 `form:"min_order_cents,optional" json:"min_order_cents,optional"`
	Starts_at int64 `form:"starts_at,optional" json:"starts_at,optional"`
	Ends_at int64 `form:"ends_at,optional" json:"ends_at,optional"`
}
```


3. response definition



```golang
type CouponSummary struct {
	Id uint64 
	Code string 
	Name string 
	Description string 
	Status string 
	Discount_type string 
	Discount_value int64 
	Currency string 
	Max_redemptions int 
	Max_redemptions_per_user int 
	Min_order_cents int64 
	Starts_at *int64 
	Ends_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Update coupon"

1. route definition

- Url: /api/v1/admin/coupons/:id
- Method: PATCH
- Request: `AdminUpdateCouponRequest`
- Response: `CouponSummary`

2. request definition



```golang
type AdminUpdateCouponRequest struct {
	Id uint64 `path:"id"`
	Name string `form:"name,optional" json:"name,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Discount_type string `form:"discount_type,optional" json:"discount_type,optional"`
	Discount_value int64 `form:"discount_value,optional" json:"discount_value,optional"`
	Currency string `form:"currency,optional" json:"currency,optional"`
	Max_redemptions int `form:"max_redemptions,optional" json:"max_redemptions,optional"`
	Max_redemptions_per_user int `form:"max_redemptions_per_user,optional" json:"max_redemptions_per_user,optional"`
	Min_order_cents int64 `form:"min_order_cents,optional" json:"min_order_cents,optional"`
	Starts_at int64 `form:"starts_at,optional" json:"starts_at,optional"`
	Ends_at int64 `form:"ends_at,optional" json:"ends_at,optional"`
}
```


3. response definition



```golang
type CouponSummary struct {
	Id uint64 
	Code string 
	Name string 
	Description string 
	Status string 
	Discount_type string 
	Discount_value int64 
	Currency string 
	Max_redemptions int 
	Max_redemptions_per_user int 
	Min_order_cents int64 
	Starts_at *int64 
	Ends_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

### 4. "Delete coupon"

1. route definition

- Url: /api/v1/admin/coupons/:id
- Method: DELETE
- Request: `AdminDeleteCouponRequest`
- Response: `-`

2. request definition



```golang
type AdminDeleteCouponRequest struct {
	Id uint64 `path:"id"`
}
```


3. response definition


