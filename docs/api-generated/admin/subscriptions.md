### 1. "List subscriptions"

1. route definition

- Url: /api/v1/admin/subscriptions
- Method: GET
- Request: `AdminListSubscriptionsRequest`
- Response: `AdminSubscriptionListResponse`

2. request definition



```golang
type AdminListSubscriptionsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status int `form:"status,optional" json:"status,optional"`
	User_id uint64 `form:"user_id,optional" json:"user_id,optional"`
	Plan_name string `form:"plan_name,optional" json:"plan_name,optional"`
	Plan_id uint64 `form:"plan_id,optional" json:"plan_id,optional"`
	Template_id uint64 `form:"template_id,optional" json:"template_id,optional"`
}
```


3. response definition



```golang
type AdminSubscriptionListResponse struct {
	Subscriptions []AdminSubscriptionSummary 
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

### 2. "Create subscription"

1. route definition

- Url: /api/v1/admin/subscriptions
- Method: POST
- Request: `AdminCreateSubscriptionRequest`
- Response: `AdminSubscriptionResponse`

2. request definition



```golang
type AdminCreateSubscriptionRequest struct {
	User_id uint64 
	Name string 
	Plan_name string `form:"plan_name,optional" json:"plan_name,optional"`
	Plan_id uint64 
	Status int `form:"status,optional" json:"status,optional"`
	Template_id uint64 
	Available_template_ids []uint64 `form:"available_template_ids,optional" json:"available_template_ids,optional"`
	Token string `form:"token,optional" json:"token,optional"`
	Expires_at int64 
	Traffic_total_bytes int64 
	Traffic_used_bytes int64 `form:"traffic_used_bytes,optional" json:"traffic_used_bytes,optional"`
	Devices_limit int 
}
```


3. response definition



```golang
type AdminSubscriptionResponse struct {
	Subscription AdminSubscriptionSummary 
}

type AdminSubscriptionSummary struct {
	Id uint64 
	User AdminSubscriptionUserSummary 
	Name string 
	Plan_name string 
	Plan_id uint64 
	Plan_snapshot map[string]interface{} 
	Status int 
	Template_id uint64 
	Available_template_ids []uint64 
	Token string 
	Expires_at int64 
	Traffic_total_bytes int64 
	Traffic_used_bytes int64 
	Devices_limit int 
	Last_refreshed_at int64 
	Created_at int64 
	Updated_at int64 
}

type AdminSubscriptionUserSummary struct {
}
```

### 3. "Get subscription detail"

1. route definition

- Url: /api/v1/admin/subscriptions/:id
- Method: GET
- Request: `AdminGetSubscriptionRequest`
- Response: `AdminSubscriptionResponse`

2. request definition



```golang
type AdminGetSubscriptionRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type AdminSubscriptionResponse struct {
	Subscription AdminSubscriptionSummary 
}

type AdminSubscriptionSummary struct {
	Id uint64 
	User AdminSubscriptionUserSummary 
	Name string 
	Plan_name string 
	Plan_id uint64 
	Plan_snapshot map[string]interface{} 
	Status int 
	Template_id uint64 
	Available_template_ids []uint64 
	Token string 
	Expires_at int64 
	Traffic_total_bytes int64 
	Traffic_used_bytes int64 
	Devices_limit int 
	Last_refreshed_at int64 
	Created_at int64 
	Updated_at int64 
}

type AdminSubscriptionUserSummary struct {
}
```

### 4. "Update subscription"

1. route definition

- Url: /api/v1/admin/subscriptions/:id
- Method: PATCH
- Request: `AdminUpdateSubscriptionRequest`
- Response: `AdminSubscriptionResponse`

2. request definition



```golang
type AdminUpdateSubscriptionRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Plan_name string `form:"plan_name,optional" json:"plan_name,optional"`
	Plan_id uint64 `form:"plan_id,optional" json:"plan_id,optional"`
	Status int `form:"status,optional" json:"status,optional"`
	Template_id uint64 `form:"template_id,optional" json:"template_id,optional"`
	Available_template_ids []uint64 `form:"available_template_ids,optional" json:"available_template_ids,optional"`
	Token string `form:"token,optional" json:"token,optional"`
	Expires_at int64 `form:"expires_at,optional" json:"expires_at,optional"`
	Traffic_total_bytes int64 `form:"traffic_total_bytes,optional" json:"traffic_total_bytes,optional"`
	Traffic_used_bytes int64 `form:"traffic_used_bytes,optional" json:"traffic_used_bytes,optional"`
	Devices_limit int `form:"devices_limit,optional" json:"devices_limit,optional"`
}
```


3. response definition



```golang
type AdminSubscriptionResponse struct {
	Subscription AdminSubscriptionSummary 
}

type AdminSubscriptionSummary struct {
	Id uint64 
	User AdminSubscriptionUserSummary 
	Name string 
	Plan_name string 
	Plan_id uint64 
	Plan_snapshot map[string]interface{} 
	Status int 
	Template_id uint64 
	Available_template_ids []uint64 
	Token string 
	Expires_at int64 
	Traffic_total_bytes int64 
	Traffic_used_bytes int64 
	Devices_limit int 
	Last_refreshed_at int64 
	Created_at int64 
	Updated_at int64 
}

type AdminSubscriptionUserSummary struct {
}
```

### 5. "Disable subscription"

1. route definition

- Url: /api/v1/admin/subscriptions/:id/disable
- Method: POST
- Request: `AdminDisableSubscriptionRequest`
- Response: `AdminSubscriptionResponse`

2. request definition



```golang
type AdminDisableSubscriptionRequest struct {
	Id uint64 
	Reason string `form:"reason,optional" json:"reason,optional"`
}
```


3. response definition



```golang
type AdminSubscriptionResponse struct {
	Subscription AdminSubscriptionSummary 
}

type AdminSubscriptionSummary struct {
	Id uint64 
	User AdminSubscriptionUserSummary 
	Name string 
	Plan_name string 
	Plan_id uint64 
	Plan_snapshot map[string]interface{} 
	Status int 
	Template_id uint64 
	Available_template_ids []uint64 
	Token string 
	Expires_at int64 
	Traffic_total_bytes int64 
	Traffic_used_bytes int64 
	Devices_limit int 
	Last_refreshed_at int64 
	Created_at int64 
	Updated_at int64 
}

type AdminSubscriptionUserSummary struct {
}
```

### 6. "Extend subscription expiry"

1. route definition

- Url: /api/v1/admin/subscriptions/:id/extend
- Method: POST
- Request: `AdminExtendSubscriptionRequest`
- Response: `AdminSubscriptionResponse`

2. request definition



```golang
type AdminExtendSubscriptionRequest struct {
	Id uint64 
	Extend_days int `form:"extend_days,optional" json:"extend_days,optional"`
	Extend_hours int `form:"extend_hours,optional" json:"extend_hours,optional"`
	Expires_at int64 `form:"expires_at,optional" json:"expires_at,optional"`
}
```


3. response definition



```golang
type AdminSubscriptionResponse struct {
	Subscription AdminSubscriptionSummary 
}

type AdminSubscriptionSummary struct {
	Id uint64 
	User AdminSubscriptionUserSummary 
	Name string 
	Plan_name string 
	Plan_id uint64 
	Plan_snapshot map[string]interface{} 
	Status int 
	Template_id uint64 
	Available_template_ids []uint64 
	Token string 
	Expires_at int64 
	Traffic_total_bytes int64 
	Traffic_used_bytes int64 
	Devices_limit int 
	Last_refreshed_at int64 
	Created_at int64 
	Updated_at int64 
}

type AdminSubscriptionUserSummary struct {
}
```

