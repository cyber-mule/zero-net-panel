### 1. "List user subscriptions"

1. route definition

- Url: /api/v1/user/subscriptions
- Method: GET
- Request: `UserListSubscriptionsRequest`
- Response: `UserSubscriptionListResponse`

2. request definition



```golang
type UserListSubscriptionsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status int `form:"status,optional" json:"status,optional"`
}
```


3. response definition



```golang
type UserSubscriptionListResponse struct {
	Subscriptions []UserSubscriptionSummary 
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

### 2. "Preview user subscription"

1. route definition

- Url: /api/v1/user/subscriptions/:id/preview
- Method: GET
- Request: `UserSubscriptionPreviewRequest`
- Response: `UserSubscriptionPreviewResponse`

2. request definition



```golang
type UserSubscriptionPreviewRequest struct {
	Id uint64 
	Template_id uint64 `form:"template_id,optional" json:"template_id,optional"`
}
```


3. response definition



```golang
type UserSubscriptionPreviewResponse struct {
	Subscription_id uint64 
	Template_id uint64 
	Content string 
	Content_type string 
	Etag string 
	Generated_at int64 
}
```

### 3. "Update user subscription template"

1. route definition

- Url: /api/v1/user/subscriptions/:id/template
- Method: POST
- Request: `UserUpdateSubscriptionTemplateRequest`
- Response: `UserUpdateSubscriptionTemplateResponse`

2. request definition



```golang
type UserUpdateSubscriptionTemplateRequest struct {
	Id uint64 
	Template_id uint64 
}
```


3. response definition



```golang
type UserUpdateSubscriptionTemplateResponse struct {
	Subscription_id uint64 
	Template_id uint64 
	Updated_at int64 
}
```

### 4. "Subscription traffic usage"

1. route definition

- Url: /api/v1/user/subscriptions/:id/traffic
- Method: GET
- Request: `UserSubscriptionTrafficRequest`
- Response: `UserSubscriptionTrafficResponse`

2. request definition



```golang
type UserSubscriptionTrafficRequest struct {
	Id uint64 
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Node_id uint64 `form:"node_id,optional" json:"node_id,optional"`
	Binding_id uint64 `form:"binding_id,optional" json:"binding_id,optional"`
	From int64 `form:"from,optional" json:"from,optional"`
	To int64 `form:"to,optional" json:"to,optional"`
}
```


3. response definition



```golang
type UserSubscriptionTrafficResponse struct {
	Summary UserSubscriptionTrafficSummary 
	Records []UserTrafficUsageRecord 
	Pagination PaginationMeta 
}

type UserSubscriptionTrafficSummary struct {
	Raw_bytes int64 
	Charged_bytes int64 
}

type PaginationMeta struct {
	Page int 
	Per_page int 
	Total_count int64 
	Has_next bool 
	Has_prev bool 
}
```

