### 1. "List billing orders"

1. route definition

- Url: /api/v1/admin/orders
- Method: GET
- Request: `AdminListOrdersRequest`
- Response: `AdminOrderListResponse`

2. request definition



```golang
type AdminListOrdersRequest struct {
	Page int 
	Per_page int 
	Status int 
	Payment_method string 
	Payment_status int 
	Number string 
	Sort string 
	Direction string 
	User_id uint64 
}
```


3. response definition



```golang
type AdminOrderListResponse struct {
	Orders []AdminOrderDetail 
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

### 2. "Get order detail"

1. route definition

- Url: /api/v1/admin/orders/:id
- Method: GET
- Request: `AdminGetOrderRequest`
- Response: `AdminOrderResponse`

2. request definition



```golang
type AdminGetOrderRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type AdminOrderResponse struct {
	Order AdminOrderDetail 
}

type AdminOrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status int 
	Payment_status int 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 3. "Cancel an order"

1. route definition

- Url: /api/v1/admin/orders/:id/cancel
- Method: POST
- Request: `AdminCancelOrderRequest`
- Response: `AdminOrderResponse`

2. request definition



```golang
type AdminCancelOrderRequest struct {
	Id uint64 
	Reason string `form:"reason,optional" json:"reason,optional"`
	Cancelled_at int64 `form:"cancelled_at,optional" json:"cancelled_at,optional"`
}
```


3. response definition



```golang
type AdminOrderResponse struct {
	Order AdminOrderDetail 
}

type AdminOrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status int 
	Payment_status int 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 4. "Manually mark an order as paid"

1. route definition

- Url: /api/v1/admin/orders/:id/pay
- Method: POST
- Request: `AdminMarkOrderPaidRequest`
- Response: `AdminOrderResponse`

2. request definition



```golang
type AdminMarkOrderPaidRequest struct {
	Id uint64 
	Payment_method string `form:"payment_method,optional" json:"payment_method,optional"`
	Paid_at int64 `form:"paid_at,optional" json:"paid_at,optional"`
	Note string `form:"note,optional" json:"note,optional"`
	Reference string `form:"reference,optional" json:"reference,optional"`
	Charge_balance bool `form:"charge_balance,optional" json:"charge_balance,optional"`
}
```


3. response definition



```golang
type AdminOrderResponse struct {
	Order AdminOrderDetail 
}

type AdminOrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status int 
	Payment_status int 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 5. "Refund an order"

1. route definition

- Url: /api/v1/admin/orders/:id/refund
- Method: POST
- Request: `AdminRefundOrderRequest`
- Response: `AdminOrderResponse`

2. request definition



```golang
type AdminRefundOrderRequest struct {
	Id uint64 
	Amount_cents int64 
	Reason string `form:"reason,optional" json:"reason,optional"`
	Metadata map[string]interface{} `form:"metadata,optional" json:"metadata,optional"`
	Refund_at int64 `form:"refund_at,optional" json:"refund_at,optional"`
	Credit_balance bool `form:"credit_balance,optional" json:"credit_balance,optional"`
}
```


3. response definition



```golang
type AdminOrderResponse struct {
	Order AdminOrderDetail 
}

type AdminOrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status int 
	Payment_status int 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 6. "Process external payment callback"

1. route definition

- Url: /api/v1/admin/orders/payments/callback
- Method: POST
- Request: `AdminPaymentCallbackRequest`
- Response: `AdminOrderResponse`

2. request definition



```golang
type AdminPaymentCallbackRequest struct {
	Order_id uint64 
	Payment_id uint64 
	Status int 
	Reference string `form:"reference,optional" json:"reference,optional"`
	Failure_code string `form:"failure_code,optional" json:"failure_code,optional"`
	Failure_message string `form:"failure_message,optional" json:"failure_message,optional"`
	Paid_at int64 `form:"paid_at,optional" json:"paid_at,optional"`
}
```


3. response definition



```golang
type AdminOrderResponse struct {
	Order AdminOrderDetail 
}

type AdminOrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status int 
	Payment_status int 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 7. "Reconcile an external payment"

1. route definition

- Url: /api/v1/admin/orders/payments/reconcile
- Method: POST
- Request: `AdminReconcilePaymentRequest`
- Response: `AdminOrderResponse`

2. request definition



```golang
type AdminReconcilePaymentRequest struct {
	Order_id uint64 
	Payment_id uint64 
}
```


3. response definition



```golang
type AdminOrderResponse struct {
	Order AdminOrderDetail 
}

type AdminOrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status int 
	Payment_status int 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 8. "Process external payment callback without admin prefix"

1. route definition

- Url: /api/v1/payments/callback
- Method: POST
- Request: `AdminPaymentCallbackRequest`
- Response: `AdminOrderResponse`

2. request definition



```golang
type AdminPaymentCallbackRequest struct {
	Order_id uint64 
	Payment_id uint64 
	Status int 
	Reference string `form:"reference,optional" json:"reference,optional"`
	Failure_code string `form:"failure_code,optional" json:"failure_code,optional"`
	Failure_message string `form:"failure_message,optional" json:"failure_message,optional"`
	Paid_at int64 `form:"paid_at,optional" json:"paid_at,optional"`
}
```


3. response definition



```golang
type AdminOrderResponse struct {
	Order AdminOrderDetail 
}

type AdminOrderDetail struct {
	Id uint64 
	Number string 
	User_id uint64 
	Status int 
	Payment_status int 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

