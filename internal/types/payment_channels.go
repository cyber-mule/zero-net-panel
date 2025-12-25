package types

// PaymentChannelSummary 支付通道摘要。
type PaymentChannelSummary struct {
	ID        uint64         `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Provider  string         `json:"provider"`
	Enabled   bool           `json:"enabled"`
	SortOrder int            `json:"sort_order"`
	Config    map[string]any `json:"config"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
}

// AdminPaymentChannelListResponse 管理端支付通道列表响应。
type AdminPaymentChannelListResponse struct {
	Channels   []PaymentChannelSummary `json:"channels"`
	Pagination PaginationMeta          `json:"pagination"`
}

// AdminListPaymentChannelsRequest 管理端支付通道列表请求。
type AdminListPaymentChannelsRequest struct {
	Page      int    `form:"page"`
	PerPage   int    `form:"per_page"`
	Query     string `form:"q"`
	Provider  string `form:"provider"`
	Enabled   *bool  `form:"enabled"`
	Sort      string `form:"sort"`
	Direction string `form:"direction"`
}

// AdminGetPaymentChannelRequest 管理端查询支付通道请求。
type AdminGetPaymentChannelRequest struct {
	ID uint64 `path:"id"`
}

// AdminCreatePaymentChannelRequest 管理端创建支付通道请求。
type AdminCreatePaymentChannelRequest struct {
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Provider  string         `json:"provider"`
	Enabled   bool           `json:"enabled"`
	SortOrder int            `json:"sort_order"`
	Config    map[string]any `json:"config"`
}

// AdminUpdatePaymentChannelRequest 管理端更新支付通道请求。
type AdminUpdatePaymentChannelRequest struct {
	ID        uint64         `path:"id"`
	Name      *string        `json:"name"`
	Code      *string        `json:"code"`
	Provider  *string        `json:"provider"`
	Enabled   *bool          `json:"enabled"`
	SortOrder *int           `json:"sort_order"`
	Config    map[string]any `json:"config"`
}
