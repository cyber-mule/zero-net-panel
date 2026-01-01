### 1. "List payment channels"

1. route definition

- Url: /api/v1/user/payment-channels
- Method: GET
- Request: `UserPaymentChannelListRequest`
- Response: `UserPaymentChannelListResponse`

2. request definition



```golang
type UserPaymentChannelListRequest struct {
	Provider string `form:"provider,optional" json:"provider,optional"`
}
```


3. response definition



```golang
type UserPaymentChannelListResponse struct {
	Channels []UserPaymentChannelSummary 
}
```

