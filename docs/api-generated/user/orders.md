### 1. "Create order from plan"

1. route definition

- Url: /api/v1/user/orders
- Method: POST
- Request: `UserCreateOrderRequest`
- Response: `UserOrderResponse`

2. request definition



```golang
type UserCreateOrderRequest struct {
	Plan_id uint64 
	Billing_option_id uint64 `form:"billing_option_id,optional" json:"billing_option_id,optional"`
	Quantity int 
	Payment_method string `form:"payment_method,optional" json:"payment_method,optional"`
	Payment_channel string `form:"payment_channel,optional" json:"payment_channel,optional"`
	Payment_return_url string `form:"payment_return_url,optional" json:"payment_return_url,optional"`
	Idempotency_key string `form:"idempotency_key,optional" json:"idempotency_key,optional"`
	Coupon_code string `form:"coupon_code,optional" json:"coupon_code,optional"`
}
```


3. response definition



```golang
type UserOrderResponse struct {
	Order OrderDetail 
	Balance BalanceSnapshot 
	Transaction *BalanceTransactionSummary 
}

type OrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status string 
	Payment_status string 
	Payment_intent_id *string 
	Payment_reference *string 
	Payment_failure_code *string 
	Payment_failure_message *string 
	Total_cents int64 
	Currency string 
	Payment_method string 
	Plan_id *uint64 
	Plan_snapshot map[string]interface{} 
	Metadata map[string]interface{} 
	Paid_at *int64 
	Cancelled_at *int64 
	Refunded_cents int64 
	Refunded_at *int64 
	Created_at int64 
	Updated_at int64 
	Items []OrderItem 
	Refunds []OrderRefund 
	Payments []OrderPayment 
}

type BalanceSnapshot struct {
	User_id uint64 
	Balance_cents int64 
	Currency string 
	Updated_at int64 
}
```

### 2. "List user orders"

1. route definition

- Url: /api/v1/user/orders
- Method: GET
- Request: `UserOrderListRequest`
- Response: `UserOrderListResponse`

2. request definition



```golang
type UserOrderListRequest struct {
	Page int 
	Per_page int 
	Status string 
	Payment_method string 
	Payment_status string 
	Number string 
	Sort string 
	Direction string 
}
```


3. response definition



```golang
type UserOrderListResponse struct {
	Orders []OrderDetail 
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

### 3. "Get user order detail"

1. route definition

- Url: /api/v1/user/orders/:id
- Method: GET
- Request: `UserGetOrderRequest`
- Response: `UserOrderResponse`

2. request definition



```golang
type UserGetOrderRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type UserOrderResponse struct {
	Order OrderDetail 
	Balance BalanceSnapshot 
	Transaction *BalanceTransactionSummary 
}

type OrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status string 
	Payment_status string 
	Payment_intent_id *string 
	Payment_reference *string 
	Payment_failure_code *string 
	Payment_failure_message *string 
	Total_cents int64 
	Currency string 
	Payment_method string 
	Plan_id *uint64 
	Plan_snapshot map[string]interface{} 
	Metadata map[string]interface{} 
	Paid_at *int64 
	Cancelled_at *int64 
	Refunded_cents int64 
	Refunded_at *int64 
	Created_at int64 
	Updated_at int64 
	Items []OrderItem 
	Refunds []OrderRefund 
	Payments []OrderPayment 
}

type BalanceSnapshot struct {
	User_id uint64 
	Balance_cents int64 
	Currency string 
	Updated_at int64 
}
```

### 4. "Cancel user order"

1. route definition

- Url: /api/v1/user/orders/:id/cancel
- Method: POST
- Request: `UserCancelOrderRequest`
- Response: `UserOrderResponse`

2. request definition



```golang
type UserCancelOrderRequest struct {
	Id uint64 
	Reason string 
}
```


3. response definition



```golang
type UserOrderResponse struct {
	Order OrderDetail 
	Balance BalanceSnapshot 
	Transaction *BalanceTransactionSummary 
}

type OrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status string 
	Payment_status string 
	Payment_intent_id *string 
	Payment_reference *string 
	Payment_failure_code *string 
	Payment_failure_message *string 
	Total_cents int64 
	Currency string 
	Payment_method string 
	Plan_id *uint64 
	Plan_snapshot map[string]interface{} 
	Metadata map[string]interface{} 
	Paid_at *int64 
	Cancelled_at *int64 
	Refunded_cents int64 
	Refunded_at *int64 
	Created_at int64 
	Updated_at int64 
	Items []OrderItem 
	Refunds []OrderRefund 
	Payments []OrderPayment 
}

type BalanceSnapshot struct {
	User_id uint64 
	Balance_cents int64 
	Currency string 
	Updated_at int64 
}
```

### 5. "Get user order payment status"

1. route definition

- Url: /api/v1/user/orders/:id/payment-status
- Method: GET
- Request: `UserOrderPaymentStatusRequest`
- Response: `UserOrderPaymentStatusResponse`

2. request definition



```golang
type UserOrderPaymentStatusRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type UserOrderPaymentStatusResponse struct {
	Order_id uint64 
	Status string 
	Payment_status string 
	Payment_method string 
	Payment_intent_id *string 
	Payment_reference *string 
	Payment_failure_code *string 
	Payment_failure_message *string 
	Paid_at *int64 
	Cancelled_at *int64 
	Refunded_cents int64 
	Refunded_at *int64 
	Updated_at int64 
}
```

