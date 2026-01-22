package types

// PlanBillingOptionSummary 套餐计费选项摘要。
type PlanBillingOptionSummary struct {
	ID            uint64 `json:"id"`
	PlanID        uint64 `json:"plan_id"`
	Name          string `json:"name"`
	DurationValue int    `json:"duration_value"`
	DurationUnit  string `json:"duration_unit"`
	PriceCents    int64  `json:"price_cents"`
	Currency      string `json:"currency"`
	SortOrder     int    `json:"sort_order"`
	Status        int    `json:"status"`
	Visible       bool   `json:"visible"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

// AdminListPlanBillingOptionsRequest 管理端套餐计费选项列表请求。
type AdminListPlanBillingOptionsRequest struct {
	PlanID  uint64 `path:"plan_id"`
	Status  int    `form:"status,optional" json:"status,optional"`
	Visible *bool  `form:"visible,optional" json:"visible,optional"`
}

// AdminPlanBillingOptionListResponse 管理端套餐计费选项列表响应。
type AdminPlanBillingOptionListResponse struct {
	Options []PlanBillingOptionSummary `json:"options"`
}

// AdminCreatePlanBillingOptionRequest 管理端创建套餐计费选项请求。
type AdminCreatePlanBillingOptionRequest struct {
	PlanID        uint64 `path:"plan_id"`
	Name          string `json:"name,omitempty,optional"`
	DurationValue int    `json:"duration_value"`
	DurationUnit  string `json:"duration_unit"`
	PriceCents    int64  `json:"price_cents"`
	Currency      string `json:"currency,omitempty,optional"`
	SortOrder     int    `json:"sort_order,omitempty,optional"`
	Status        int    `json:"status,omitempty,optional"`
	Visible       bool   `json:"visible,omitempty,optional"`
}

// AdminUpdatePlanBillingOptionRequest 管理端更新套餐计费选项请求。
type AdminUpdatePlanBillingOptionRequest struct {
	PlanID        uint64  `path:"plan_id"`
	OptionID      uint64  `path:"id"`
	Name          *string `json:"name,omitempty,optional"`
	DurationValue *int    `json:"duration_value,omitempty,optional"`
	DurationUnit  *string `json:"duration_unit,omitempty,optional"`
	PriceCents    *int64  `json:"price_cents,omitempty,optional"`
	Currency      *string `json:"currency,omitempty,optional"`
	SortOrder     *int    `json:"sort_order,omitempty,optional"`
	Status        *int    `json:"status,omitempty,optional"`
	Visible       *bool   `json:"visible,omitempty,optional"`
}
