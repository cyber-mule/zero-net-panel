### 1. "List subscription plans"

1. route definition

- Url: /api/v1/admin/plans
- Method: GET
- Request: `AdminListPlansRequest`
- Response: `AdminPlanListResponse`

2. request definition



```golang
type AdminListPlansRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Visible bool `form:"visible,optional" json:"visible,optional"`
}
```


3. response definition



```golang
type AdminPlanListResponse struct {
	Plans []PlanSummary 
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

### 2. "Create subscription plan"

1. route definition

- Url: /api/v1/admin/plans
- Method: POST
- Request: `AdminCreatePlanRequest`
- Response: `PlanSummary`

2. request definition



```golang
type AdminCreatePlanRequest struct {
	Name string 
	Slug string `form:"slug,optional" json:"slug,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Features []string `form:"features,optional" json:"features,optional"`
	Binding_ids []uint64 `form:"binding_ids,optional" json:"binding_ids,optional"`
	Price_cents int64 
	Currency string 
	Duration_days int 
	Traffic_limit_bytes int64 `form:"traffic_limit_bytes,optional" json:"traffic_limit_bytes,optional"`
	Traffic_multipliers map[string]float64 `form:"traffic_multipliers,optional" json:"traffic_multipliers,optional"`
	Devices_limit int `form:"devices_limit,optional" json:"devices_limit,optional"`
	Sort_order int `form:"sort_order,optional" json:"sort_order,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Visible bool `form:"visible,optional" json:"visible,optional"`
}
```


3. response definition



```golang
type PlanSummary struct {
	Id uint64 
	Name string 
	Slug string 
	Description string 
	Tags []string 
	Features []string 
	Binding_ids []uint64 
	Billing_options []PlanBillingOptionSummary 
	Price_cents int64 
	Currency string 
	Duration_days int 
	Traffic_limit_bytes int64 
	Traffic_multipliers map[string]float64 
	Devices_limit int 
	Sort_order int 
	Status string 
	Visible bool 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Update subscription plan"

1. route definition

- Url: /api/v1/admin/plans/:id
- Method: PATCH
- Request: `AdminUpdatePlanRequest`
- Response: `PlanSummary`

2. request definition



```golang
type AdminUpdatePlanRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Slug string `form:"slug,optional" json:"slug,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Features []string `form:"features,optional" json:"features,optional"`
	Binding_ids []uint64 `form:"binding_ids,optional" json:"binding_ids,optional"`
	Price_cents int64 `form:"price_cents,optional" json:"price_cents,optional"`
	Currency string `form:"currency,optional" json:"currency,optional"`
	Duration_days int `form:"duration_days,optional" json:"duration_days,optional"`
	Traffic_limit_bytes int64 `form:"traffic_limit_bytes,optional" json:"traffic_limit_bytes,optional"`
	Traffic_multipliers map[string]float64 `form:"traffic_multipliers,optional" json:"traffic_multipliers,optional"`
	Devices_limit int `form:"devices_limit,optional" json:"devices_limit,optional"`
	Sort_order int `form:"sort_order,optional" json:"sort_order,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Visible bool `form:"visible,optional" json:"visible,optional"`
}
```


3. response definition



```golang
type PlanSummary struct {
	Id uint64 
	Name string 
	Slug string 
	Description string 
	Tags []string 
	Features []string 
	Binding_ids []uint64 
	Billing_options []PlanBillingOptionSummary 
	Price_cents int64 
	Currency string 
	Duration_days int 
	Traffic_limit_bytes int64 
	Traffic_multipliers map[string]float64 
	Devices_limit int 
	Sort_order int 
	Status string 
	Visible bool 
	Created_at int64 
	Updated_at int64 
}
```

