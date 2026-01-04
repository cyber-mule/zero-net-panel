### 1. "List protocol entries"

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

### 2. "Create protocol entry"

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

### 3. "Update protocol entry"

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

### 4. "Delete protocol entry"

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


