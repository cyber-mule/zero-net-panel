### 1. "List protocol configs"

1. route definition

- Url: /api/v1/admin/protocol-configs
- Method: GET
- Request: `AdminListProtocolConfigsRequest`
- Response: `AdminProtocolConfigListResponse`

2. request definition



```golang
type AdminListProtocolConfigsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Status string `form:"status,optional" json:"status,optional"`
}
```


3. response definition



```golang
type AdminProtocolConfigListResponse struct {
	Configs []ProtocolConfigSummary 
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

### 2. "Create protocol config"

1. route definition

- Url: /api/v1/admin/protocol-configs
- Method: POST
- Request: `AdminCreateProtocolConfigRequest`
- Response: `ProtocolConfigSummary`

2. request definition



```golang
type AdminCreateProtocolConfigRequest struct {
	Name string 
	Protocol string 
	Status string `form:"status,optional" json:"status,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Profile map[string]interface{} `form:"profile,optional" json:"profile,optional"`
}
```


3. response definition



```golang
type ProtocolConfigSummary struct {
	Id uint64 
	Name string 
	Protocol string 
	Status string 
	Tags []string 
	Description string 
	Profile map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Update protocol config"

1. route definition

- Url: /api/v1/admin/protocol-configs/:id
- Method: PATCH
- Request: `AdminUpdateProtocolConfigRequest`
- Response: `ProtocolConfigSummary`

2. request definition



```golang
type AdminUpdateProtocolConfigRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Tags []string `form:"tags,optional" json:"tags,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Profile map[string]interface{} `form:"profile,optional" json:"profile,optional"`
}
```


3. response definition



```golang
type ProtocolConfigSummary struct {
	Id uint64 
	Name string 
	Protocol string 
	Status string 
	Tags []string 
	Description string 
	Profile map[string]interface{} 
	Created_at int64 
	Updated_at int64 
}
```

### 4. "Delete protocol config"

1. route definition

- Url: /api/v1/admin/protocol-configs/:id
- Method: DELETE
- Request: `AdminDeleteProtocolConfigRequest`
- Response: `-`

2. request definition



```golang
type AdminDeleteProtocolConfigRequest struct {
	Id uint64 
}
```


3. response definition


