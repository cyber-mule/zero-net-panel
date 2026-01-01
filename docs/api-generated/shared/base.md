### 1. "Service health check"

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

