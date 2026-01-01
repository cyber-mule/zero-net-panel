package types

// AdminListSubscriptionsRequest filters admin subscription list.
type AdminListSubscriptionsRequest struct {
	Page       int    `form:"page,optional" json:"page,optional"`
	PerPage    int    `form:"per_page,optional" json:"per_page,optional"`
	Query      string `form:"q,optional" json:"q,optional"`
	Status     string `form:"status,optional" json:"status,optional"`
	UserID     uint64 `form:"user_id,optional" json:"user_id,optional"`
	PlanName   string `form:"plan_name,optional" json:"plan_name,optional"`
	PlanID     uint64 `form:"plan_id,optional" json:"plan_id,optional"`
	TemplateID uint64 `form:"template_id,optional" json:"template_id,optional"`
}

// AdminSubscriptionUserSummary returns user info for subscription.
type AdminSubscriptionUserSummary struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

// AdminSubscriptionSummary describes subscription data for admin.
type AdminSubscriptionSummary struct {
	ID                   uint64                       `json:"id"`
	User                 AdminSubscriptionUserSummary `json:"user"`
	Name                 string                       `json:"name"`
	PlanName             string                       `json:"plan_name"`
	PlanID               uint64                       `json:"plan_id"`
	PlanSnapshot         map[string]any               `json:"plan_snapshot"`
	Status               string                       `json:"status"`
	TemplateID           uint64                       `json:"template_id"`
	AvailableTemplateIDs []uint64                     `json:"available_template_ids"`
	Token                string                       `json:"token"`
	ExpiresAt            int64                        `json:"expires_at"`
	TrafficTotalBytes    int64                        `json:"traffic_total_bytes"`
	TrafficUsedBytes     int64                        `json:"traffic_used_bytes"`
	DevicesLimit         int                          `json:"devices_limit"`
	LastRefreshedAt      int64                        `json:"last_refreshed_at"`
	CreatedAt            int64                        `json:"created_at"`
	UpdatedAt            int64                        `json:"updated_at"`
}

// AdminSubscriptionListResponse returns paginated subscriptions.
type AdminSubscriptionListResponse struct {
	Subscriptions []AdminSubscriptionSummary `json:"subscriptions"`
	Pagination    PaginationMeta             `json:"pagination"`
}

// AdminGetSubscriptionRequest fetches subscription by ID.
type AdminGetSubscriptionRequest struct {
	SubscriptionID uint64 `path:"id"`
}

// AdminSubscriptionResponse returns a subscription.
type AdminSubscriptionResponse struct {
	Subscription AdminSubscriptionSummary `json:"subscription"`
}

// AdminCreateSubscriptionRequest creates a subscription.
type AdminCreateSubscriptionRequest struct {
	UserID               uint64   `json:"user_id"`
	Name                 string   `json:"name"`
	PlanName             string   `json:"plan_name,omitempty,optional"`
	PlanID               uint64   `json:"plan_id"`
	Status               *string  `json:"status,omitempty,optional"`
	TemplateID           uint64   `json:"template_id"`
	AvailableTemplateIDs []uint64 `json:"available_template_ids,omitempty,optional"`
	Token                *string  `json:"token,omitempty,optional"`
	ExpiresAt            int64    `json:"expires_at"`
	TrafficTotalBytes    int64    `json:"traffic_total_bytes"`
	TrafficUsedBytes     *int64   `json:"traffic_used_bytes,omitempty,optional"`
	DevicesLimit         int      `json:"devices_limit"`
}

// AdminUpdateSubscriptionRequest updates subscription fields.
type AdminUpdateSubscriptionRequest struct {
	SubscriptionID       uint64    `path:"id"`
	Name                 *string   `json:"name,omitempty,optional"`
	PlanName             *string   `json:"plan_name,omitempty,optional"`
	PlanID               *uint64   `json:"plan_id,omitempty,optional"`
	Status               *string   `json:"status,omitempty,optional"`
	TemplateID           *uint64   `json:"template_id,omitempty,optional"`
	AvailableTemplateIDs *[]uint64 `json:"available_template_ids,omitempty,optional"`
	Token                *string   `json:"token,omitempty,optional"`
	ExpiresAt            *int64    `json:"expires_at,omitempty,optional"`
	TrafficTotalBytes    *int64    `json:"traffic_total_bytes,omitempty,optional"`
	TrafficUsedBytes     *int64    `json:"traffic_used_bytes,omitempty,optional"`
	DevicesLimit         *int      `json:"devices_limit,omitempty,optional"`
}

// AdminDisableSubscriptionRequest disables a subscription.
type AdminDisableSubscriptionRequest struct {
	SubscriptionID uint64  `path:"id"`
	Reason         *string `json:"reason,omitempty,optional"`
}

// AdminExtendSubscriptionRequest extends subscription expiry.
type AdminExtendSubscriptionRequest struct {
	SubscriptionID uint64 `path:"id"`
	ExtendDays     int    `json:"extend_days,omitempty,optional"`
	ExtendHours    int    `json:"extend_hours,omitempty,optional"`
	ExpiresAt      *int64 `json:"expires_at,omitempty,optional"`
}
