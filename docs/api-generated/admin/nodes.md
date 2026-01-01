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
	Protocols []string `form:"protocols,optional" json:"protocols,optional"`
	Capacity_mbps int `form:"capacity_mbps,optional" json:"capacity_mbps,optional"`
	Description string `form:"description,optional" json:"description,optional"`
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
	Protocols []string 
	Capacity_mbps int 
	Description string 
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
	Protocols []string `form:"protocols,optional" json:"protocols,optional"`
	Capacity_mbps int `form:"capacity_mbps,optional" json:"capacity_mbps,optional"`
	Description string `form:"description,optional" json:"description,optional"`
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
	Protocols []string 
	Capacity_mbps int 
	Description string 
	Last_synced_at int64 
	Updated_at int64 
}
```

### 4. "Disable node"

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
	Protocols []string 
	Capacity_mbps int 
	Description string 
	Last_synced_at int64 
	Updated_at int64 
}
```

### 5. "Get node kernel endpoints"

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

### 6. "Upsert node kernel endpoint"

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

### 7. "Sync node kernel configuration"

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

