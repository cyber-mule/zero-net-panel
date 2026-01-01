### 1. "List available plans"

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

