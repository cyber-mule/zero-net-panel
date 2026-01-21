### 1. "List edge nodes"

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

### 2. "Create node"

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

### 3. "Update node"

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

### 4. "Delete node (soft delete)"

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


### 5. "Disable node"

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

### 6. "Get node kernel endpoints"

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

### 7. "Upsert node kernel endpoint"

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

### 8. "Sync node kernel configuration"

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

### 9. "Sync node status"

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

