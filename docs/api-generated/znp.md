### 1. "List announcements"

1. route definition

- Url: /api/v1/admin/announcements
- Method: GET
- Request: `AdminListAnnouncementsRequest`
- Response: `AdminAnnouncementListResponse`

2. request definition



```golang
type AdminListAnnouncementsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Category string `form:"category,optional" json:"category,optional"`
	Audience string `form:"audience,optional" json:"audience,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
}
```


3. response definition



```golang
type AdminAnnouncementListResponse struct {
	Announcements []AnnouncementSummary 
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

### 2. "Create announcement"

1. route definition

- Url: /api/v1/admin/announcements
- Method: POST
- Request: `AdminCreateAnnouncementRequest`
- Response: `AnnouncementSummary`

2. request definition



```golang
type AdminCreateAnnouncementRequest struct {
	Title string 
	Content string 
	Category string `form:"category,optional" json:"category,optional"`
	Audience string `form:"audience,optional" json:"audience,optional"`
	Is_pinned bool `form:"is_pinned,optional" json:"is_pinned,optional"`
	Priority int `form:"priority,optional" json:"priority,optional"`
	Created_by string `form:"created_by,optional" json:"created_by,optional"`
}
```


3. response definition



```golang
type AnnouncementSummary struct {
	Id uint64 
	Title string 
	Content string 
	Category string 
	Status string 
	Audience string 
	Is_pinned bool 
	Priority int 
	Visible_from int64 
	Visible_to int64 `form:"visible_to,optional" json:"visible_to,optional"`
	Published_at int64 `form:"published_at,optional" json:"published_at,optional"`
	Published_by string 
	Created_by string 
	Updated_by string 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Publish announcement"

1. route definition

- Url: /api/v1/admin/announcements/:id/publish
- Method: POST
- Request: `AdminPublishAnnouncementRequest`
- Response: `AnnouncementSummary`

2. request definition



```golang
type AdminPublishAnnouncementRequest struct {
	Id uint64 
	Visible_to int64 `form:"visible_to,optional" json:"visible_to,optional"`
	Operator string `form:"operator,optional" json:"operator,optional"`
}
```


3. response definition



```golang
type AnnouncementSummary struct {
	Id uint64 
	Title string 
	Content string 
	Category string 
	Status string 
	Audience string 
	Is_pinned bool 
	Priority int 
	Visible_from int64 
	Visible_to int64 `form:"visible_to,optional" json:"visible_to,optional"`
	Published_at int64 `form:"published_at,optional" json:"published_at,optional"`
	Published_by string 
	Created_by string 
	Updated_by string 
	Created_at int64 
	Updated_at int64 
}
```

### 4. "List audit logs"

1. route definition

- Url: /api/v1/admin/audit-logs
- Method: GET
- Request: `AdminAuditLogListRequest`
- Response: `AdminAuditLogListResponse`

2. request definition



```golang
type AdminAuditLogListRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Actor_id *uint64 `form:"actor_id,optional" json:"actor_id,optional"`
	Action string `form:"action,optional" json:"action,optional"`
	Resource_type string `form:"resource_type,optional" json:"resource_type,optional"`
	Resource_id string `form:"resource_id,optional" json:"resource_id,optional"`
	Since int64 `form:"since,optional" json:"since,optional"`
	Until int64 `form:"until,optional" json:"until,optional"`
}
```


3. response definition



```golang
type AdminAuditLogListResponse struct {
	Logs []AuditLogSummary 
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

### 5. "Export audit logs"

1. route definition

- Url: /api/v1/admin/audit-logs/export
- Method: GET
- Request: `AdminAuditLogExportRequest`
- Response: `AdminAuditLogExportResponse`

2. request definition



```golang
type AdminAuditLogExportRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Actor_id *uint64 `form:"actor_id,optional" json:"actor_id,optional"`
	Action string `form:"action,optional" json:"action,optional"`
	Resource_type string `form:"resource_type,optional" json:"resource_type,optional"`
	Resource_id string `form:"resource_id,optional" json:"resource_id,optional"`
	Since int64 `form:"since,optional" json:"since,optional"`
	Until int64 `form:"until,optional" json:"until,optional"`
	Format string `form:"format,optional" json:"format,optional"`
}
```


3. response definition



```golang
type AdminAuditLogExportResponse struct {
	Logs []AuditLogSummary 
	Total_count int64 
	Exported_at int64 
}
```

### 6. "List coupons"

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

### 7. "Create coupon"

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

### 8. "Update coupon"

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

### 9. "Delete coupon"

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


### 10. "List admin console modules"

1. route definition

- Url: /api/v1/admin/dashboard
- Method: GET
- Request: `-`
- Response: `AdminDashboardResponse`

2. request definition



3. response definition



```golang
type AdminDashboardResponse struct {
	Modules []AdminModule 
}
```

### 11. "List edge nodes"

1. route definition

- Url: /api/v1/admin/nodes
- Method: GET
- Request: `AdminListNodesRequest`
- Response: `AdminNodeListResponse`

2. request definition



```golang
type AdminListNodesRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
}
```


3. response definition



```golang
type AdminNodeListResponse struct {
	Nodes []NodeSummary 
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

### 12. "Create node"

1. route definition

- Url: /api/v1/admin/nodes
- Method: POST
- Request: `AdminCreateNodeRequest`
- Response: `AdminNodeResponse`

2. request definition



```golang
type AdminCreateNodeRequest struct {
	Name string 
	Region string `form:"region,optional" json:"region,optional"`
	Country string `form:"country,optional" json:"country,optional"`
	Isp string `form:"isp,optional" json:"isp,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Capacity_mbps int `form:"capacity_mbps,optional" json:"capacity_mbps,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Access_address string `form:"access_address,optional" json:"access_address,optional"`
	Control_endpoint string `form:"control_endpoint" json:"control_endpoint"`
	Control_access_key string `form:"control_access_key,optional" json:"control_access_key,optional"`
	Control_secret_key string `form:"control_secret_key,optional" json:"control_secret_key,optional"`
	Ak string `form:"ak,optional" json:"ak,optional"`
	Sk string `form:"sk,optional" json:"sk,optional"`
	Control_token string `form:"control_token,optional" json:"control_token,optional"`
	Kernel_default_protocol string `form:"kernel_default_protocol,optional" json:"kernel_default_protocol,optional"`
	Kernel_http_timeout_seconds int `form:"kernel_http_timeout_seconds,optional" json:"kernel_http_timeout_seconds,optional"`
	Kernel_status_poll_interval_seconds int `form:"kernel_status_poll_interval_seconds,optional" json:"kernel_status_poll_interval_seconds,optional"`
	Kernel_status_poll_backoff_enabled bool `form:"kernel_status_poll_backoff_enabled,optional" json:"kernel_status_poll_backoff_enabled,optional"`
	Kernel_status_poll_backoff_max_interval_seconds int `form:"kernel_status_poll_backoff_max_interval_seconds,optional" json:"kernel_status_poll_backoff_max_interval_seconds,optional"`
	Kernel_status_poll_backoff_multiplier float64 `form:"kernel_status_poll_backoff_multiplier,optional" json:"kernel_status_poll_backoff_multiplier,optional"`
	Kernel_status_poll_backoff_jitter float64 `form:"kernel_status_poll_backoff_jitter,optional" json:"kernel_status_poll_backoff_jitter,optional"`
	Kernel_offline_probe_max_interval_seconds int `form:"kernel_offline_probe_max_interval_seconds,optional" json:"kernel_offline_probe_max_interval_seconds,optional"`
	Status_sync_enabled bool `form:"status_sync_enabled,optional" json:"status_sync_enabled,optional"`
}
```


3. response definition



```golang
type AdminNodeResponse struct {
	Node NodeSummary 
}

type NodeSummary struct {
	Id uint64 
	Name string 
	Region string 
	Country string 
	Isp string 
	Status string 
	Tags []string 
	Capacity_mbps int 
	Description string 
	Access_address string 
	Control_endpoint string 
	Kernel_default_protocol string 
	Kernel_http_timeout_seconds int 
	Kernel_status_poll_interval_seconds int 
	Kernel_status_poll_backoff_enabled bool 
	Kernel_status_poll_backoff_max_interval_seconds int 
	Kernel_status_poll_backoff_multiplier float64 
	Kernel_status_poll_backoff_jitter float64 
	Kernel_offline_probe_max_interval_seconds int 
	Status_sync_enabled bool 
	Last_synced_at int64 
	Updated_at int64 
}
```

### 13. "Update node"

1. route definition

- Url: /api/v1/admin/nodes/:id
- Method: PATCH
- Request: `AdminUpdateNodeRequest`
- Response: `AdminNodeResponse`

2. request definition



```golang
type AdminUpdateNodeRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Region string `form:"region,optional" json:"region,optional"`
	Country string `form:"country,optional" json:"country,optional"`
	Isp string `form:"isp,optional" json:"isp,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Capacity_mbps int `form:"capacity_mbps,optional" json:"capacity_mbps,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Access_address string `form:"access_address,optional" json:"access_address,optional"`
	Control_endpoint string `form:"control_endpoint,optional" json:"control_endpoint,optional"`
	Control_access_key string `form:"control_access_key,optional" json:"control_access_key,optional"`
	Control_secret_key string `form:"control_secret_key,optional" json:"control_secret_key,optional"`
	Ak string `form:"ak,optional" json:"ak,optional"`
	Sk string `form:"sk,optional" json:"sk,optional"`
	Control_token string `form:"control_token,optional" json:"control_token,optional"`
	Kernel_default_protocol string `form:"kernel_default_protocol,optional" json:"kernel_default_protocol,optional"`
	Kernel_http_timeout_seconds int `form:"kernel_http_timeout_seconds,optional" json:"kernel_http_timeout_seconds,optional"`
	Kernel_status_poll_interval_seconds int `form:"kernel_status_poll_interval_seconds,optional" json:"kernel_status_poll_interval_seconds,optional"`
	Kernel_status_poll_backoff_enabled bool `form:"kernel_status_poll_backoff_enabled,optional" json:"kernel_status_poll_backoff_enabled,optional"`
	Kernel_status_poll_backoff_max_interval_seconds int `form:"kernel_status_poll_backoff_max_interval_seconds,optional" json:"kernel_status_poll_backoff_max_interval_seconds,optional"`
	Kernel_status_poll_backoff_multiplier float64 `form:"kernel_status_poll_backoff_multiplier,optional" json:"kernel_status_poll_backoff_multiplier,optional"`
	Kernel_status_poll_backoff_jitter float64 `form:"kernel_status_poll_backoff_jitter,optional" json:"kernel_status_poll_backoff_jitter,optional"`
	Kernel_offline_probe_max_interval_seconds int `form:"kernel_offline_probe_max_interval_seconds,optional" json:"kernel_offline_probe_max_interval_seconds,optional"`
	Status_sync_enabled bool `form:"status_sync_enabled,optional" json:"status_sync_enabled,optional"`
}
```


3. response definition



```golang
type AdminNodeResponse struct {
	Node NodeSummary 
}

type NodeSummary struct {
	Id uint64 
	Name string 
	Region string 
	Country string 
	Isp string 
	Status string 
	Tags []string 
	Capacity_mbps int 
	Description string 
	Access_address string 
	Control_endpoint string 
	Kernel_default_protocol string 
	Kernel_http_timeout_seconds int 
	Kernel_status_poll_interval_seconds int 
	Kernel_status_poll_backoff_enabled bool 
	Kernel_status_poll_backoff_max_interval_seconds int 
	Kernel_status_poll_backoff_multiplier float64 
	Kernel_status_poll_backoff_jitter float64 
	Kernel_offline_probe_max_interval_seconds int 
	Status_sync_enabled bool 
	Last_synced_at int64 
	Updated_at int64 
}
```

### 14. "Delete node (soft delete)"

1. route definition

- Url: /api/v1/admin/nodes/:id
- Method: DELETE
- Request: `AdminDeleteNodeRequest`
- Response: `-`

2. request definition



```golang
type AdminDeleteNodeRequest struct {
	Id uint64 
}
```


3. response definition


### 15. "Disable node"

1. route definition

- Url: /api/v1/admin/nodes/:id/disable
- Method: POST
- Request: `AdminDisableNodeRequest`
- Response: `AdminNodeResponse`

2. request definition



```golang
type AdminDisableNodeRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type AdminNodeResponse struct {
	Node NodeSummary 
}

type NodeSummary struct {
	Id uint64 
	Name string 
	Region string 
	Country string 
	Isp string 
	Status string 
	Tags []string 
	Capacity_mbps int 
	Description string 
	Access_address string 
	Control_endpoint string 
	Kernel_default_protocol string 
	Kernel_http_timeout_seconds int 
	Kernel_status_poll_interval_seconds int 
	Kernel_status_poll_backoff_enabled bool 
	Kernel_status_poll_backoff_max_interval_seconds int 
	Kernel_status_poll_backoff_multiplier float64 
	Kernel_status_poll_backoff_jitter float64 
	Kernel_offline_probe_max_interval_seconds int 
	Status_sync_enabled bool 
	Last_synced_at int64 
	Updated_at int64 
}
```

### 16. "Get node kernel endpoints"

1. route definition

- Url: /api/v1/admin/nodes/:id/kernels
- Method: GET
- Request: `AdminNodeKernelPath`
- Response: `AdminNodeKernelResponse`

2. request definition



```golang
type AdminNodeKernelPath struct {
	Id uint64 
}
```


3. response definition



```golang
type AdminNodeKernelResponse struct {
	Node_id uint64 
	Kernels []NodeKernelSummary 
}
```

### 17. "Upsert node kernel endpoint"

1. route definition

- Url: /api/v1/admin/nodes/:id/kernels
- Method: POST
- Request: `AdminUpsertNodeKernelRequest`
- Response: `AdminNodeKernelUpsertResponse`

2. request definition



```golang
type AdminUpsertNodeKernelRequest struct {
	Id uint64 
	Protocol string 
	Endpoint string 
	Revision string `form:"revision,optional" json:"revision,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Config map[string]interface{} `form:"config,optional" json:"config,optional"`
	Last_synced_at int64 `form:"last_synced_at,optional" json:"last_synced_at,optional"`
}
```


3. response definition



```golang
type AdminNodeKernelUpsertResponse struct {
	Node_id uint64 
	Kernel NodeKernelSummary 
}

type NodeKernelSummary struct {
	Protocol string 
	Endpoint string 
	Revision string 
	Status string 
	Config map[string]interface{} 
	Last_synced_at int64 
}
```

### 18. "Sync node kernel configuration"

1. route definition

- Url: /api/v1/admin/nodes/:id/kernels/sync
- Method: POST
- Request: `AdminSyncNodeKernelRequest`
- Response: `AdminSyncNodeKernelResponse`

2. request definition



```golang
type AdminSyncNodeKernelRequest struct {
	Id uint64 
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
}
```


3. response definition



```golang
type AdminSyncNodeKernelResponse struct {
	Node_id uint64 
	Protocol string 
	Revision string 
	Synced_at int64 
	Message string 
}
```

### 19. "Sync node status"

1. route definition

- Url: /api/v1/admin/nodes/status/sync
- Method: POST
- Request: `AdminSyncNodeStatusRequest`
- Response: `AdminSyncNodeStatusResponse`

2. request definition



```golang
type AdminSyncNodeStatusRequest struct {
	Node_ids []uint64 `form:"node_ids,optional" json:"node_ids,optional"`
}
```


3. response definition



```golang
type AdminSyncNodeStatusResponse struct {
	Results []NodeStatusSyncResult 
}
```

### 20. "List billing orders"

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
	Status string 
	Payment_method string 
	Payment_status string 
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

### 21. "Get order detail"

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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 22. "Cancel an order"

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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 23. "Manually mark an order as paid"

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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 24. "Refund an order"

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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 25. "Process external payment callback"

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
	Status string 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 26. "Reconcile an external payment"

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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 27. "Process external payment callback without admin prefix"

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
	Status string 
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
	User OrderUserSummary 
}

type OrderUserSummary struct {
}
```

### 28. "List payment channels"

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

### 29. "Create payment channel"

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

### 30. "Get payment channel"

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

### 31. "Update payment channel"

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

### 32. "List plan billing options"

1. route definition

- Url: /api/v1/admin/plans/:plan_id/billing-options
- Method: GET
- Request: `AdminListPlanBillingOptionsRequest`
- Response: `AdminPlanBillingOptionListResponse`

2. request definition



```golang
type AdminListPlanBillingOptionsRequest struct {
	Plan_id uint64 `path:"plan_id"`
	Status string `form:"status,optional" json:"status,optional"`
	Visible bool `form:"visible,optional" json:"visible,optional"`
}
```


3. response definition



```golang
type AdminPlanBillingOptionListResponse struct {
	Options []PlanBillingOptionSummary 
}
```

### 33. "Create plan billing option"

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
	Status string `form:"status,optional" json:"status,optional"`
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
	Status string 
	Visible bool 
	Created_at int64 
	Updated_at int64 
}
```

### 34. "Update plan billing option"

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
	Status string `form:"status,optional" json:"status,optional"`
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
	Status string 
	Visible bool 
	Created_at int64 
	Updated_at int64 
}
```

### 35. "List subscription plans"

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

### 36. "Create subscription plan"

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

### 37. "Update subscription plan"

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

### 38. "List protocol bindings"

1. route definition

- Url: /api/v1/admin/protocol-bindings
- Method: GET
- Request: `AdminListProtocolBindingsRequest`
- Response: `AdminProtocolBindingListResponse`

2. request definition



```golang
type AdminListProtocolBindingsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Node_id uint64 `form:"node_id,optional" json:"node_id,optional"`
}
```


3. response definition



```golang
type AdminProtocolBindingListResponse struct {
	Bindings []ProtocolBindingSummary 
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

### 39. "Create protocol binding"

1. route definition

- Url: /api/v1/admin/protocol-bindings
- Method: POST
- Request: `AdminCreateProtocolBindingRequest`
- Response: `ProtocolBindingSummary`

2. request definition



```golang
type AdminCreateProtocolBindingRequest struct {
	Name string `form:"name,optional" json:"name,optional"`
	Node_id uint64 
	Protocol string 
	Profile map[string]interface{} 
	Role string 
	Listen string `form:"listen,optional" json:"listen,optional"`
	Connect string `form:"connect,optional" json:"connect,optional"`
	Access_port int `form:"access_port,optional" json:"access_port,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Kernel_id string `form:"kernel_id" json:"kernel_id"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Metadata map[string]interface{} `form:"metadata,optional" json:"metadata,optional"`
}
```


3. response definition



```golang
type ProtocolBindingSummary struct {
	Id uint64 
	Name string 
	Node_id uint64 
	Node_name string 
	Protocol string 
	Role string 
	Listen string 
	Connect string 
	Access_port int 
	Status string 
	Kernel_id string 
	Sync_status string 
	Health_status string 
	Last_synced_at int64 
	Last_heartbeat_at int64 
	Last_sync_error string 
	Tags []string 
	Description string 
	Profile map[string]interface{} 
	Metadata map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 40. "Update protocol binding"

1. route definition

- Url: /api/v1/admin/protocol-bindings/:id
- Method: PATCH
- Request: `AdminUpdateProtocolBindingRequest`
- Response: `ProtocolBindingSummary`

2. request definition



```golang
type AdminUpdateProtocolBindingRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Node_id uint64 `form:"node_id,optional" json:"node_id,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Role string `form:"role,optional" json:"role,optional"`
	Listen string `form:"listen,optional" json:"listen,optional"`
	Connect string `form:"connect,optional" json:"connect,optional"`
	Access_port int `form:"access_port,optional" json:"access_port,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Kernel_id string `form:"kernel_id,optional" json:"kernel_id,optional"`
	Sync_status string `form:"sync_status,optional" json:"sync_status,optional"`
	Health_status string `form:"health_status,optional" json:"health_status,optional"`
	Last_synced_at int64 `form:"last_synced_at,optional" json:"last_synced_at,optional"`
	Last_heartbeat_at int64 `form:"last_heartbeat_at,optional" json:"last_heartbeat_at,optional"`
	Last_sync_error string `form:"last_sync_error,optional" json:"last_sync_error,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Profile map[string]interface{} `form:"profile,optional" json:"profile,optional"`
	Metadata map[string]interface{} `form:"metadata,optional" json:"metadata,optional"`
}
```


3. response definition



```golang
type ProtocolBindingSummary struct {
	Id uint64 
	Name string 
	Node_id uint64 
	Node_name string 
	Protocol string 
	Role string 
	Listen string 
	Connect string 
	Access_port int 
	Status string 
	Kernel_id string 
	Sync_status string 
	Health_status string 
	Last_synced_at int64 
	Last_heartbeat_at int64 
	Last_sync_error string 
	Tags []string 
	Description string 
	Profile map[string]interface{} 
	Metadata map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 41. "Delete protocol binding"

1. route definition

- Url: /api/v1/admin/protocol-bindings/:id
- Method: DELETE
- Request: `AdminDeleteProtocolBindingRequest`
- Response: `-`

2. request definition



```golang
type AdminDeleteProtocolBindingRequest struct {
	Id uint64 
}
```


3. response definition


### 42. "Sync protocol binding"

1. route definition

- Url: /api/v1/admin/protocol-bindings/:id/sync
- Method: POST
- Request: `AdminSyncProtocolBindingRequest`
- Response: `ProtocolBindingSyncResult`

2. request definition



```golang
type AdminSyncProtocolBindingRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type ProtocolBindingSyncResult struct {
	Binding_id uint64 
	Status string 
	Message string 
	Synced_at int64 
}
```

### 43. "Sync protocol binding status"

1. route definition

- Url: /api/v1/admin/protocol-bindings/status/sync
- Method: POST
- Request: `AdminSyncProtocolBindingStatusRequest`
- Response: `AdminSyncProtocolBindingStatusResponse`

2. request definition



```golang
type AdminSyncProtocolBindingStatusRequest struct {
	Node_ids []uint64 `form:"node_ids,optional" json:"node_ids,optional"`
}
```


3. response definition



```golang
type AdminSyncProtocolBindingStatusResponse struct {
	Results []ProtocolBindingStatusSyncResult 
}
```

### 44. "Sync protocol bindings"

1. route definition

- Url: /api/v1/admin/protocol-bindings/sync
- Method: POST
- Request: `AdminSyncProtocolBindingsRequest`
- Response: `AdminSyncProtocolBindingsResponse`

2. request definition



```golang
type AdminSyncProtocolBindingsRequest struct {
	Binding_ids []uint64 `form:"binding_ids,optional" json:"binding_ids,optional"`
	Node_ids []uint64 `form:"node_ids,optional" json:"node_ids,optional"`
}
```


3. response definition



```golang
type AdminSyncProtocolBindingsResponse struct {
	Results []ProtocolBindingSyncResult 
}
```

### 45. "List protocol entries"

1. route definition

- Url: /api/v1/admin/protocol-entries
- Method: GET
- Request: `AdminListProtocolEntriesRequest`
- Response: `AdminProtocolEntryListResponse`

2. request definition



```golang
type AdminListProtocolEntriesRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Binding_id uint64 `form:"binding_id,optional" json:"binding_id,optional"`
}
```


3. response definition



```golang
type AdminProtocolEntryListResponse struct {
	Entries []ProtocolEntrySummary 
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

### 46. "Create protocol entry"

1. route definition

- Url: /api/v1/admin/protocol-entries
- Method: POST
- Request: `AdminCreateProtocolEntryRequest`
- Response: `ProtocolEntrySummary`

2. request definition



```golang
type AdminCreateProtocolEntryRequest struct {
	Name string `form:"name,optional" json:"name,optional"`
	Binding_id uint64 
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Entry_address string `form:"entry_address" json:"entry_address"`
	Entry_port int `form:"entry_port" json:"entry_port"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Profile map[string]interface{} `form:"profile,optional" json:"profile,optional"`
}
```


3. response definition



```golang
type ProtocolEntrySummary struct {
	Id uint64 
	Name string 
	Binding_id uint64 
	Binding_name string 
	Node_id uint64 
	Node_name string 
	Protocol string 
	Status string 
	Binding_status string 
	Health_status string 
	Entry_address string 
	Entry_port int 
	Tags []string 
	Description string 
	Profile map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 47. "Update protocol entry"

1. route definition

- Url: /api/v1/admin/protocol-entries/:id
- Method: PATCH
- Request: `AdminUpdateProtocolEntryRequest`
- Response: `ProtocolEntrySummary`

2. request definition



```golang
type AdminUpdateProtocolEntryRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Binding_id uint64 `form:"binding_id,optional" json:"binding_id,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Entry_address string `form:"entry_address,optional" json:"entry_address,optional"`
	Entry_port int `form:"entry_port,optional" json:"entry_port,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Profile map[string]interface{} `form:"profile,optional" json:"profile,optional"`
}
```


3. response definition



```golang
type ProtocolEntrySummary struct {
	Id uint64 
	Name string 
	Binding_id uint64 
	Binding_name string 
	Node_id uint64 
	Node_name string 
	Protocol string 
	Status string 
	Binding_status string 
	Health_status string 
	Entry_address string 
	Entry_port int 
	Tags []string 
	Description string 
	Profile map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 48. "Delete protocol entry"

1. route definition

- Url: /api/v1/admin/protocol-entries/:id
- Method: DELETE
- Request: `AdminDeleteProtocolEntryRequest`
- Response: `-`

2. request definition



```golang
type AdminDeleteProtocolEntryRequest struct {
	Id uint64 
}
```


3. response definition


### 49. "Get third-party security settings"

1. route definition

- Url: /api/v1/admin/security-settings
- Method: GET
- Request: `-`
- Response: `AdminSecuritySettingResponse`

2. request definition



3. response definition



```golang
type AdminSecuritySettingResponse struct {
	Setting SecuritySetting 
}

type SecuritySetting struct {
	Id uint64 
	Third_party_api_enabled bool 
	Api_key string 
	Api_secret string 
	Encryption_algorithm string 
	Nonce_ttl_seconds int 
	Created_at int64 
	Updated_at int64 
}
```

### 50. "Update third-party security settings"

1. route definition

- Url: /api/v1/admin/security-settings
- Method: PATCH
- Request: `AdminUpdateSecuritySettingRequest`
- Response: `AdminSecuritySettingResponse`

2. request definition



```golang
type AdminUpdateSecuritySettingRequest struct {
	Third_party_api_enabled bool `form:"third_party_api_enabled,optional" json:"third_party_api_enabled,optional"`
	Api_key string `form:"api_key,optional" json:"api_key,optional"`
	Api_secret string `form:"api_secret,optional" json:"api_secret,optional"`
	Encryption_algorithm string `form:"encryption_algorithm,optional" json:"encryption_algorithm,optional"`
	Nonce_ttl_seconds int `form:"nonce_ttl_seconds,optional" json:"nonce_ttl_seconds,optional"`
}
```


3. response definition



```golang
type AdminSecuritySettingResponse struct {
	Setting SecuritySetting 
}

type SecuritySetting struct {
	Id uint64 
	Third_party_api_enabled bool 
	Api_key string 
	Api_secret string 
	Encryption_algorithm string 
	Nonce_ttl_seconds int 
	Created_at int64 
	Updated_at int64 
}
```

### 51. "Get site settings"

1. route definition

- Url: /api/v1/admin/site-settings
- Method: GET
- Request: `-`
- Response: `AdminSiteSettingResponse`

2. request definition



3. response definition



```golang
type AdminSiteSettingResponse struct {
	Setting SiteSetting 
}

type SiteSetting struct {
	Id uint64 
	Name string 
	Logo_url string 
	Access_domain string 
	Created_at int64 
	Updated_at int64 
}
```

### 52. "Update site settings"

1. route definition

- Url: /api/v1/admin/site-settings
- Method: PATCH
- Request: `AdminUpdateSiteSettingRequest`
- Response: `AdminSiteSettingResponse`

2. request definition



```golang
type AdminUpdateSiteSettingRequest struct {
	Name string `form:"name,optional" json:"name,optional"`
	Logo_url string `form:"logo_url,optional" json:"logo_url,optional"`
	Access_domain string `form:"access_domain,optional" json:"access_domain,optional"`
}
```


3. response definition



```golang
type AdminSiteSettingResponse struct {
	Setting SiteSetting 
}

type SiteSetting struct {
	Id uint64 
	Name string 
	Logo_url string 
	Access_domain string 
	Created_at int64 
	Updated_at int64 
}
```

### 53. "List subscriptions"

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
	Status string `form:"status,optional" json:"status,optional"`
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

### 54. "Create subscription"

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
	Status string `form:"status,optional" json:"status,optional"`
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
	Status string 
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

### 55. "Get subscription detail"

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
	Status string 
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

### 56. "Update subscription"

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
	Status string `form:"status,optional" json:"status,optional"`
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
	Status string 
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

### 57. "Disable subscription"

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
	Status string 
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

### 58. "Extend subscription expiry"

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
	Status string 
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

### 59. "List subscription templates"

1. route definition

- Url: /api/v1/admin/subscription-templates
- Method: GET
- Request: `AdminListSubscriptionTemplatesRequest`
- Response: `AdminSubscriptionTemplateListResponse`

2. request definition



```golang
type AdminListSubscriptionTemplatesRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Client_type string `form:"client_type,optional" json:"client_type,optional"`
	Format string `form:"format,optional" json:"format,optional"`
	Include_drafts bool `form:"include_drafts,optional" json:"include_drafts,optional"`
}
```


3. response definition



```golang
type AdminSubscriptionTemplateListResponse struct {
	Templates []SubscriptionTemplateSummary 
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

### 60. "Create subscription template"

1. route definition

- Url: /api/v1/admin/subscription-templates
- Method: POST
- Request: `AdminCreateSubscriptionTemplateRequest`
- Response: `SubscriptionTemplateSummary`

2. request definition



```golang
type AdminCreateSubscriptionTemplateRequest struct {
	Name string 
	Description string `form:"description,optional" json:"description,optional"`
	Client_type string 
	Format string `form:"format,optional" json:"format,optional"`
	Content string 
	Variables map[string]TemplateVariable `form:"variables,optional" json:"variables,optional"`
	Is_default bool `form:"is_default,optional" json:"is_default,optional"`
}
```


3. response definition



```golang
type SubscriptionTemplateSummary struct {
	Id uint64 
	Name string 
	Description string 
	Client_type string 
	Format string 
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable 
	Is_default bool 
	Version uint32 
	Updated_at int64 
	Published_at int64 
	Last_published_by string 
}
```

### 61. "Update subscription template"

1. route definition

- Url: /api/v1/admin/subscription-templates/:id
- Method: PATCH
- Request: `AdminUpdateSubscriptionTemplateRequest`
- Response: `SubscriptionTemplateSummary`

2. request definition



```golang
type AdminUpdateSubscriptionTemplateRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Format string `form:"format,optional" json:"format,optional"`
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable `form:"variables,optional" json:"variables,optional"`
	Is_default bool `form:"is_default,optional" json:"is_default,optional"`
}
```


3. response definition



```golang
type SubscriptionTemplateSummary struct {
	Id uint64 
	Name string 
	Description string 
	Client_type string 
	Format string 
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable 
	Is_default bool 
	Version uint32 
	Updated_at int64 
	Published_at int64 
	Last_published_by string 
}
```

### 62. "List template publish history"

1. route definition

- Url: /api/v1/admin/subscription-templates/:id/history
- Method: GET
- Request: `AdminSubscriptionTemplateHistoryRequest`
- Response: `AdminSubscriptionTemplateHistoryResponse`

2. request definition



```golang
type AdminSubscriptionTemplateHistoryRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type AdminSubscriptionTemplateHistoryResponse struct {
	Template_id uint64 
	History []SubscriptionTemplateHistoryEntry 
}
```

### 63. "Publish subscription template"

1. route definition

- Url: /api/v1/admin/subscription-templates/:id/publish
- Method: POST
- Request: `AdminPublishSubscriptionTemplateRequest`
- Response: `AdminPublishSubscriptionTemplateResponse`

2. request definition



```golang
type AdminPublishSubscriptionTemplateRequest struct {
	Id uint64 
	Changelog string `form:"changelog,optional" json:"changelog,optional"`
	Operator string `form:"operator,optional" json:"operator,optional"`
}
```


3. response definition



```golang
type AdminPublishSubscriptionTemplateResponse struct {
	Template SubscriptionTemplateSummary 
	History SubscriptionTemplateHistoryEntry 
}

type SubscriptionTemplateSummary struct {
	Id uint64 
	Name string 
	Description string 
	Client_type string 
	Format string 
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable 
	Is_default bool 
	Version uint32 
	Updated_at int64 
	Published_at int64 
	Last_published_by string 
}

type SubscriptionTemplateHistoryEntry struct {
	Version uint32 
	Changelog string 
	Published_at int64 
	Published_by string 
	Variables map[string]TemplateVariable 
}
```

### 64. "List users with filters"

1. route definition

- Url: /api/v1/admin/users
- Method: GET
- Request: `AdminListUsersRequest`
- Response: `AdminUserListResponse`

2. request definition



```golang
type AdminListUsersRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Role string `form:"role,optional" json:"role,optional"`
}
```


3. response definition



```golang
type AdminUserListResponse struct {
	Users []AdminUserSummary 
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

### 65. "Create a new user"

1. route definition

- Url: /api/v1/admin/users
- Method: POST
- Request: `AdminCreateUserRequest`
- Response: `AdminUserResponse`

2. request definition



```golang
type AdminCreateUserRequest struct {
	Email string 
	Password string 
	Display_name string `form:"display_name,optional" json:"display_name,optional"`
	Roles []string `form:"roles,optional" json:"roles,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Email_verified bool `form:"email_verified,optional" json:"email_verified,optional"`
}
```


3. response definition



```golang
type AdminUserResponse struct {
	User AdminUserSummary 
}

type AdminUserSummary struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Status string 
	Email_verified_at int64 `form:"email_verified_at,optional" json:"email_verified_at,optional"`
	Failed_login_attempts int 
	Locked_until int64 `form:"locked_until,optional" json:"locked_until,optional"`
	Last_login_at int64 `form:"last_login_at,optional" json:"last_login_at,optional"`
	Created_at int64 
	Updated_at int64 
}
```

### 66. "Rotate user credential"

1. route definition

- Url: /api/v1/admin/users/:id/credentials/rotate
- Method: POST
- Request: `AdminRotateUserCredentialRequest`
- Response: `AdminRotateUserCredentialResponse`

2. request definition



```golang
type AdminRotateUserCredentialRequest struct {
	Id uint64 `path:"id"`
}
```


3. response definition



```golang
type AdminRotateUserCredentialResponse struct {
	User_id uint64 
	Credential CredentialSummary 
}

type CredentialSummary struct {
	Version int 
	Status string 
	Issued_at int64 
	Deprecated_at *int64 
	Revoked_at *int64 
	Last_seen_at *int64 
}
```

### 67. "Force user logout"

1. route definition

- Url: /api/v1/admin/users/:id/force-logout
- Method: POST
- Request: `AdminForceLogoutRequest`
- Response: `AdminForceLogoutResponse`

2. request definition



```golang
type AdminForceLogoutRequest struct {
	Id uint64 `path:"id"`
}
```


3. response definition



```golang
type AdminForceLogoutResponse struct {
	Message string 
}
```

### 68. "Reset user password"

1. route definition

- Url: /api/v1/admin/users/:id/reset-password
- Method: POST
- Request: `AdminResetUserPasswordRequest`
- Response: `AdminResetUserPasswordResponse`

2. request definition



```golang
type AdminResetUserPasswordRequest struct {
	Id uint64 `path:"id"`
	Password string 
}
```


3. response definition



```golang
type AdminResetUserPasswordResponse struct {
	Message string 
}
```

### 69. "Update user roles"

1. route definition

- Url: /api/v1/admin/users/:id/roles
- Method: PATCH
- Request: `AdminUpdateUserRolesRequest`
- Response: `AdminUserResponse`

2. request definition



```golang
type AdminUpdateUserRolesRequest struct {
	Id uint64 `path:"id"`
	Roles []string 
}
```


3. response definition



```golang
type AdminUserResponse struct {
	User AdminUserSummary 
}

type AdminUserSummary struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Status string 
	Email_verified_at int64 `form:"email_verified_at,optional" json:"email_verified_at,optional"`
	Failed_login_attempts int 
	Locked_until int64 `form:"locked_until,optional" json:"locked_until,optional"`
	Last_login_at int64 `form:"last_login_at,optional" json:"last_login_at,optional"`
	Created_at int64 
	Updated_at int64 
}
```

### 70. "Update user status (active/disabled)"

1. route definition

- Url: /api/v1/admin/users/:id/status
- Method: PATCH
- Request: `AdminUpdateUserStatusRequest`
- Response: `AdminUserResponse`

2. request definition



```golang
type AdminUpdateUserStatusRequest struct {
	Id uint64 `path:"id"`
	Status string 
}
```


3. response definition



```golang
type AdminUserResponse struct {
	User AdminUserSummary 
}

type AdminUserSummary struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Status string 
	Email_verified_at int64 `form:"email_verified_at,optional" json:"email_verified_at,optional"`
	Failed_login_attempts int 
	Locked_until int64 `form:"locked_until,optional" json:"locked_until,optional"`
	Last_login_at int64 `form:"last_login_at,optional" json:"last_login_at,optional"`
	Created_at int64 
	Updated_at int64 
}
```

### 71. "Request password reset code"

1. route definition

- Url: /api/v1/auth/forgot
- Method: POST
- Request: `AuthForgotPasswordRequest`
- Response: `AuthForgotPasswordResponse`

2. request definition



```golang
type AuthForgotPasswordRequest struct {
	Email string 
}
```


3. response definition



```golang
type AuthForgotPasswordResponse struct {
	Message string 
}
```

### 72. "Authenticate user and issue token"

1. route definition

- Url: /api/v1/auth/login
- Method: POST
- Request: `AuthLoginRequest`
- Response: `AuthLoginResponse`

2. request definition



```golang
type AuthLoginRequest struct {
	Email string 
	Password string 
}
```


3. response definition



```golang
type AuthLoginResponse struct {
	Access_token string 
	Refresh_token string 
	Token_type string 
	Expires_in int64 
	Refresh_expires_in int64 
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

### 73. "Refresh access token"

1. route definition

- Url: /api/v1/auth/refresh
- Method: POST
- Request: `AuthRefreshRequest`
- Response: `AuthRefreshResponse`

2. request definition



```golang
type AuthRefreshRequest struct {
	Refresh_token string 
}
```


3. response definition



```golang
type AuthRefreshResponse struct {
	Access_token string 
	Refresh_token string 
	Token_type string 
	Expires_in int64 
	Refresh_expires_in int64 
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

### 74. "Register new user account"

1. route definition

- Url: /api/v1/auth/register
- Method: POST
- Request: `AuthRegisterRequest`
- Response: `AuthRegisterResponse`

2. request definition



```golang
type AuthRegisterRequest struct {
	Email string 
	Password string 
	Display_name string `form:"display_name,optional" json:"display_name,optional"`
	Invite_code string `form:"invite_code,optional" json:"invite_code,optional"`
}
```


3. response definition



```golang
type AuthRegisterResponse struct {
	Requires_verification bool 
	Access_token string `form:"access_token,optional" json:"access_token,optional"`
	Refresh_token string `form:"refresh_token,optional" json:"refresh_token,optional"`
	Token_type string `form:"token_type,optional" json:"token_type,optional"`
	Expires_in int64 `form:"expires_in,optional" json:"expires_in,optional"`
	Refresh_expires_in int64 `form:"refresh_expires_in,optional" json:"refresh_expires_in,optional"`
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

### 75. "Reset password using verification code"

1. route definition

- Url: /api/v1/auth/reset
- Method: POST
- Request: `AuthResetPasswordRequest`
- Response: `AuthResetPasswordResponse`

2. request definition



```golang
type AuthResetPasswordRequest struct {
	Email string 
	Code string 
	Password string 
}
```


3. response definition



```golang
type AuthResetPasswordResponse struct {
	Message string 
}
```

### 76. "Verify user email with code"

1. route definition

- Url: /api/v1/auth/verify
- Method: POST
- Request: `AuthVerifyRequest`
- Response: `AuthVerifyResponse`

2. request definition



```golang
type AuthVerifyRequest struct {
	Email string 
	Code string 
}
```


3. response definition



```golang
type AuthVerifyResponse struct {
	Access_token string 
	Refresh_token string 
	Token_type string 
	Expires_in int64 
	Refresh_expires_in int64 
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

### 77. "Service health check"

1. route definition

- Url: /api/v1/ping
- Method: GET
- Request: `-`
- Response: `PingResponse`

2. request definition



3. response definition



```golang
type PingResponse struct {
	Status string 
	Service string 
	Version string 
	Site_name string 
	Logo_url string 
	Timestamp int64 
}
```

### 78. "Get user balance"

1. route definition

- Url: /api/v1/user/account/balance
- Method: GET
- Request: `UserBalanceRequest`
- Response: `UserBalanceResponse`

2. request definition



```golang
type UserBalanceRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Entry_type string `form:"entry_type,optional" json:"entry_type,optional"`
}
```


3. response definition



```golang
type UserBalanceResponse struct {
	User_id uint64 
	Balance_cents int64 
	Currency string 
	Updated_at int64 
	Transactions []BalanceTransactionSummary 
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

### 79. "Rotate credential"

1. route definition

- Url: /api/v1/user/account/credentials/rotate
- Method: POST
- Request: `UserRotateCredentialRequest`
- Response: `UserRotateCredentialResponse`

2. request definition



```golang
type UserRotateCredentialRequest struct {
}
```


3. response definition



```golang
type UserRotateCredentialResponse struct {
	Credential CredentialSummary 
}

type CredentialSummary struct {
	Version int 
	Status string 
	Issued_at int64 
	Deprecated_at *int64 
	Revoked_at *int64 
	Last_seen_at *int64 
}
```

### 80. "Change email"

1. route definition

- Url: /api/v1/user/account/email
- Method: POST
- Request: `UserChangeEmailRequest`
- Response: `UserChangeEmailResponse`

2. request definition



```golang
type UserChangeEmailRequest struct {
	Email string 
	Code string 
	Password string 
}
```


3. response definition



```golang
type UserChangeEmailResponse struct {
	Profile UserProfile 
}

type UserProfile struct {
	Id uint64 
	Email string 
	Display_name string 
	Status string 
	Email_verified_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

### 81. "Send email change code"

1. route definition

- Url: /api/v1/user/account/email/code
- Method: POST
- Request: `UserEmailChangeCodeRequest`
- Response: `UserEmailChangeCodeResponse`

2. request definition



```golang
type UserEmailChangeCodeRequest struct {
	Email string 
}
```


3. response definition



```golang
type UserEmailChangeCodeResponse struct {
	Message string 
}
```

### 82. "Change password"

1. route definition

- Url: /api/v1/user/account/password
- Method: POST
- Request: `UserChangePasswordRequest`
- Response: `UserChangePasswordResponse`

2. request definition



```golang
type UserChangePasswordRequest struct {
	Current_password string 
	New_password string 
}
```


3. response definition



```golang
type UserChangePasswordResponse struct {
	Message string 
}
```

### 83. "Get user profile"

1. route definition

- Url: /api/v1/user/account/profile
- Method: GET
- Request: `UserProfileRequest`
- Response: `UserProfileResponse`

2. request definition



```golang
type UserProfileRequest struct {
}
```


3. response definition



```golang
type UserProfileResponse struct {
	Profile UserProfile 
}

type UserProfile struct {
	Id uint64 
	Email string 
	Display_name string 
	Status string 
	Email_verified_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

### 84. "Update user profile"

1. route definition

- Url: /api/v1/user/account/profile
- Method: PATCH
- Request: `UserUpdateProfileRequest`
- Response: `UserProfileResponse`

2. request definition



```golang
type UserUpdateProfileRequest struct {
	Display_name string `form:"display_name,optional" json:"display_name,optional"`
}
```


3. response definition



```golang
type UserProfileResponse struct {
	Profile UserProfile 
}

type UserProfile struct {
	Id uint64 
	Email string 
	Display_name string 
	Status string 
	Email_verified_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

### 85. "List active announcements"

1. route definition

- Url: /api/v1/user/announcements
- Method: GET
- Request: `UserAnnouncementListRequest`
- Response: `UserAnnouncementListResponse`

2. request definition



```golang
type UserAnnouncementListRequest struct {
	Audience string `form:"audience,optional" json:"audience,optional"`
	Limit int `form:"limit,optional" json:"limit,optional"`
}
```


3. response definition



```golang
type UserAnnouncementListResponse struct {
	Announcements []UserAnnouncementSummary 
}
```

### 86. "List node runtime status (sanitized)"

1. route definition

- Url: /api/v1/user/nodes
- Method: GET
- Request: `UserNodeStatusListRequest`
- Response: `UserNodeStatusListResponse`

2. request definition



```golang
type UserNodeStatusListRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
}
```


3. response definition



```golang
type UserNodeStatusListResponse struct {
	Nodes []UserNodeStatusSummary 
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

### 87. "Create order from plan"

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

### 88. "List user orders"

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

### 89. "Get user order detail"

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

### 90. "Cancel user order"

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

### 91. "Get user order payment status"

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

### 92. "List payment channels"

1. route definition

- Url: /api/v1/user/payment-channels
- Method: GET
- Request: `UserPaymentChannelListRequest`
- Response: `UserPaymentChannelListResponse`

2. request definition



```golang
type UserPaymentChannelListRequest struct {
	Provider string `form:"provider,optional" json:"provider,optional"`
}
```


3. response definition



```golang
type UserPaymentChannelListResponse struct {
	Channels []UserPaymentChannelSummary 
}
```

### 93. "List available plans"

1. route definition

- Url: /api/v1/user/plans
- Method: GET
- Request: `UserPlanListRequest`
- Response: `UserPlanListResponse`

2. request definition



```golang
type UserPlanListRequest struct {
	Q string `form:"q,optional" json:"q,optional"`
}
```


3. response definition



```golang
type UserPlanListResponse struct {
	Plans []UserPlanSummary 
}
```

### 94. "List user subscriptions"

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
	Status string `form:"status,optional" json:"status,optional"`
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

### 95. "Preview user subscription"

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

### 96. "Update user subscription template"

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

### 97. "Subscription traffic usage"

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

