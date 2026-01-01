### 1. "List protocol bindings"

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
	Protocol_config_id uint64 `form:"protocol_config_id,optional" json:"protocol_config_id,optional"`
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

### 2. "Create protocol binding"

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
	Protocol_config_id uint64 `form:"protocol_config_id,optional" json:"protocol_config_id,optional"`
	Role string 
	Listen string `form:"listen,optional" json:"listen,optional"`
	Connect string `form:"connect,optional" json:"connect,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Kernel_id string `form:"kernel_id,optional" json:"kernel_id,optional"`
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
	Protocol_config_id uint64 
	Protocol string 
	Role string 
	Listen string 
	Connect string 
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

### 3. "Update protocol binding"

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
	Protocol_config_id uint64 `form:"protocol_config_id,optional" json:"protocol_config_id,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Role string `form:"role,optional" json:"role,optional"`
	Listen string `form:"listen,optional" json:"listen,optional"`
	Connect string `form:"connect,optional" json:"connect,optional"`
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
	Protocol_config_id uint64 
	Protocol string 
	Role string 
	Listen string 
	Connect string 
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

### 4. "Delete protocol binding"

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


### 5. "Sync protocol binding"

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

### 6. "Sync protocol bindings"

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

