### 1. "List payment channels"

1. route definition

- Url: /api/v1/admin/payment-channels
- Method: GET
- Request: `AdminListPaymentChannelsRequest`
- Response: `AdminPaymentChannelListResponse`

2. request definition



```golang
type AdminListPaymentChannelsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Provider string `form:"provider,optional" json:"provider,optional"`
	Enabled bool `form:"enabled,optional" json:"enabled,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
}
```


3. response definition



```golang
type AdminPaymentChannelListResponse struct {
	Channels []PaymentChannelSummary 
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

### 2. "Create payment channel"

1. route definition

- Url: /api/v1/admin/payment-channels
- Method: POST
- Request: `AdminCreatePaymentChannelRequest`
- Response: `PaymentChannelSummary`

2. request definition



```golang
type AdminCreatePaymentChannelRequest struct {
	Name string 
	Code string 
	Provider string `form:"provider,optional" json:"provider,optional"`
	Enabled bool `form:"enabled,optional" json:"enabled,optional"`
	Sort_order int `form:"sort_order,optional" json:"sort_order,optional"`
	Config map[string]interface{} `form:"config,optional" json:"config,optional"`
}
```


3. response definition



```golang
type PaymentChannelSummary struct {
	Id uint64 
	Name string 
	Code string 
	Provider string 
	Enabled bool 
	Sort_order int 
	Config map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Get payment channel"

1. route definition

- Url: /api/v1/admin/payment-channels/:id
- Method: GET
- Request: `AdminGetPaymentChannelRequest`
- Response: `PaymentChannelSummary`

2. request definition



```golang
type AdminGetPaymentChannelRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type PaymentChannelSummary struct {
	Id uint64 
	Name string 
	Code string 
	Provider string 
	Enabled bool 
	Sort_order int 
	Config map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 4. "Update payment channel"

1. route definition

- Url: /api/v1/admin/payment-channels/:id
- Method: PATCH
- Request: `AdminUpdatePaymentChannelRequest`
- Response: `PaymentChannelSummary`

2. request definition



```golang
type AdminUpdatePaymentChannelRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Code string `form:"code,optional" json:"code,optional"`
	Provider string `form:"provider,optional" json:"provider,optional"`
	Enabled bool `form:"enabled,optional" json:"enabled,optional"`
	Sort_order int `form:"sort_order,optional" json:"sort_order,optional"`
	Config map[string]interface{} `form:"config,optional" json:"config,optional"`
}
```


3. response definition



```golang
type PaymentChannelSummary struct {
	Id uint64 
	Name string 
	Code string 
	Provider string 
	Enabled bool 
	Sort_order int 
	Config map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

