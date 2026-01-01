### 1. "List audit logs"

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

### 2. "Export audit logs"

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

