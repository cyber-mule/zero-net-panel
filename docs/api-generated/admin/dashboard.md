### 1. "List admin console modules"

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

