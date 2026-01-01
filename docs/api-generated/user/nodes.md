### 1. "List node runtime status (sanitized)"

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

