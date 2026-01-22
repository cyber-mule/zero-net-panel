### 1. "List plan billing options"

1. route definition

- Url: /api/v1/admin/plans/:plan_id/billing-options
- Method: GET
- Request: `AdminListPlanBillingOptionsRequest`
- Response: `AdminPlanBillingOptionListResponse`

2. request definition



```golang
type AdminListPlanBillingOptionsRequest struct {
	Plan_id uint64 `path:"plan_id"`
	Status int `form:"status,optional" json:"status,optional"`
	Visible bool `form:"visible,optional" json:"visible,optional"`
}
```


3. response definition



```golang
type AdminPlanBillingOptionListResponse struct {
	Options []PlanBillingOptionSummary 
}
```

### 2. "Create plan billing option"

1. route definition

- Url: /api/v1/admin/plans/:plan_id/billing-options
- Method: POST
- Request: `AdminCreatePlanBillingOptionRequest`
- Response: `PlanBillingOptionSummary`

2. request definition



```golang
type AdminCreatePlanBillingOptionRequest struct {
	Plan_id uint64 `path:"plan_id"`
	Name string `form:"name,optional" json:"name,optional"`
	Duration_value int 
	Duration_unit string 
	Price_cents int64 
	Currency string `form:"currency,optional" json:"currency,optional"`
	Sort_order int `form:"sort_order,optional" json:"sort_order,optional"`
	Status int `form:"status,optional" json:"status,optional"`
	Visible bool `form:"visible,optional" json:"visible,optional"`
}
```


3. response definition



```golang
type PlanBillingOptionSummary struct {
	Id uint64 
	Plan_id uint64 
	Name string 
	Duration_value int 
	Duration_unit string 
	Price_cents int64 
	Currency string 
	Sort_order int 
	Status int 
	Visible bool 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Update plan billing option"

1. route definition

- Url: /api/v1/admin/plans/:plan_id/billing-options/:id
- Method: PATCH
- Request: `AdminUpdatePlanBillingOptionRequest`
- Response: `PlanBillingOptionSummary`

2. request definition



```golang
type AdminUpdatePlanBillingOptionRequest struct {
	Plan_id uint64 `path:"plan_id"`
	Id uint64 `path:"id"`
	Name string `form:"name,optional" json:"name,optional"`
	Duration_value int `form:"duration_value,optional" json:"duration_value,optional"`
	Duration_unit string `form:"duration_unit,optional" json:"duration_unit,optional"`
	Price_cents int64 `form:"price_cents,optional" json:"price_cents,optional"`
	Currency string `form:"currency,optional" json:"currency,optional"`
	Sort_order int `form:"sort_order,optional" json:"sort_order,optional"`
	Status int `form:"status,optional" json:"status,optional"`
	Visible bool `form:"visible,optional" json:"visible,optional"`
}
```


3. response definition



```golang
type PlanBillingOptionSummary struct {
	Id uint64 
	Plan_id uint64 
	Name string 
	Duration_value int 
	Duration_unit string 
	Price_cents int64 
	Currency string 
	Sort_order int 
	Status int 
	Visible bool 
	Created_at int64 
	Updated_at int64 
}
```

