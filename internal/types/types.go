package types

// PingResponse 保留健康检查响应。
type PingResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	SiteName  string `json:"site_name"`
	LogoURL   string `json:"logo_url"`
	Timestamp int64  `json:"timestamp"`
}

// PaginationMeta 统一 GitHub 风格分页返回。
type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalCount int64 `json:"total_count"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// AuthLoginRequest 登录请求。
type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthRefreshRequest 刷新令牌请求。
type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthenticatedUser 鉴权用户信息。
type AuthenticatedUser struct {
	ID          uint64   `json:"id"`
	Email       string   `json:"email"`
	DisplayName string   `json:"display_name"`
	Roles       []string `json:"roles"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

// AuthLoginResponse 登录响应。
type AuthLoginResponse struct {
	AccessToken      string            `json:"access_token"`
	RefreshToken     string            `json:"refresh_token"`
	TokenType        string            `json:"token_type"`
	ExpiresIn        int64             `json:"expires_in"`
	RefreshExpiresIn int64             `json:"refresh_expires_in"`
	User             AuthenticatedUser `json:"user"`
}

// AuthRefreshResponse 刷新响应。
type AuthRefreshResponse struct {
	AccessToken      string            `json:"access_token"`
	RefreshToken     string            `json:"refresh_token"`
	TokenType        string            `json:"token_type"`
	ExpiresIn        int64             `json:"expires_in"`
	RefreshExpiresIn int64             `json:"refresh_expires_in"`
	User             AuthenticatedUser `json:"user"`
}

// AdminModule 管理后台模块信息。
type AdminModule struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Route       string   `json:"route"`
	Permissions []string `json:"permissions"`
}

// AdminDashboardResponse 返回管理后台模块集合。
type AdminDashboardResponse struct {
	Modules []AdminModule `json:"modules"`
}

// AuditLogSummary 审计日志摘要。
type AuditLogSummary struct {
	ID           uint64         `json:"id"`
	ActorID      *uint64        `json:"actor_id"`
	ActorEmail   string         `json:"actor_email"`
	ActorRoles   []string       `json:"actor_roles"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resource_type"`
	ResourceID   string         `json:"resource_id"`
	SourceIP     string         `json:"source_ip"`
	Metadata     map[string]any `json:"metadata"`
	CreatedAt    int64          `json:"created_at"`
}

// AdminAuditLogListRequest 审计日志列表请求。
type AdminAuditLogListRequest struct {
	Page         int     `form:"page,optional" json:"page,optional"`
	PerPage      int     `form:"per_page,optional" json:"per_page,optional"`
	ActorID      *uint64 `form:"actor_id,optional" json:"actor_id,optional"`
	Action       string  `form:"action,optional" json:"action,optional"`
	ResourceType string  `form:"resource_type,optional" json:"resource_type,optional"`
	ResourceID   string  `form:"resource_id,optional" json:"resource_id,optional"`
	Since        int64   `form:"since,optional" json:"since,optional"`
	Until        int64   `form:"until,optional" json:"until,optional"`
}

// AdminAuditLogListResponse 审计日志列表响应。
type AdminAuditLogListResponse struct {
	Logs       []AuditLogSummary `json:"logs"`
	Pagination PaginationMeta    `json:"pagination"`
}

// AdminAuditLogExportRequest 审计日志导出请求。
type AdminAuditLogExportRequest struct {
	Page         int     `form:"page,optional" json:"page,optional"`
	PerPage      int     `form:"per_page,optional" json:"per_page,optional"`
	ActorID      *uint64 `form:"actor_id,optional" json:"actor_id,optional"`
	Action       string  `form:"action,optional" json:"action,optional"`
	ResourceType string  `form:"resource_type,optional" json:"resource_type,optional"`
	ResourceID   string  `form:"resource_id,optional" json:"resource_id,optional"`
	Since        int64   `form:"since,optional" json:"since,optional"`
	Until        int64   `form:"until,optional" json:"until,optional"`
	Format       string  `form:"format,optional" json:"format,optional"`
}

// AdminAuditLogExportResponse 审计日志导出响应。
type AdminAuditLogExportResponse struct {
	Logs       []AuditLogSummary `json:"logs"`
	TotalCount int64             `json:"total_count"`
	ExportedAt int64             `json:"exported_at"`
}

// SiteSetting 站点品牌配置。
type SiteSetting struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	LogoURL   string `json:"logo_url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// AdminSiteSettingResponse 站点配置响应。
type AdminSiteSettingResponse struct {
	Setting SiteSetting `json:"setting"`
}

// AdminUpdateSiteSettingRequest 更新站点配置。
type AdminUpdateSiteSettingRequest struct {
	Name    *string `json:"name,optional"`
	LogoURL *string `json:"logo_url,optional"`
}

// AdminListNodesRequest 管理端节点列表查询参数。
type AdminListNodesRequest struct {
	Page      int    `form:"page,optional" json:"page,optional"`
	PerPage   int    `form:"per_page,optional" json:"per_page,optional"`
	Sort      string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Query     string `form:"q,optional" json:"q,optional"`
	Status    string `form:"status,optional" json:"status,optional"`
	Protocol  string `form:"protocol,optional" json:"protocol,optional"`
}

// NodeSummary 节点摘要信息。
type NodeSummary struct {
	ID           uint64   `json:"id"`
	Name         string   `json:"name"`
	Region       string   `json:"region"`
	Country      string   `json:"country"`
	ISP          string   `json:"isp"`
	Status       string   `json:"status"`
	Tags         []string `json:"tags"`
	Protocols    []string `json:"protocols"`
	CapacityMbps int      `json:"capacity_mbps"`
	Description  string   `json:"description"`
	LastSyncedAt int64    `json:"last_synced_at"`
	UpdatedAt    int64    `json:"updated_at"`
}

// AdminNodeResponse 节点详情响应。
type AdminNodeResponse struct {
	Node NodeSummary `json:"node"`
}

// AdminCreateNodeRequest 管理端创建节点请求。
type AdminCreateNodeRequest struct {
	Name         string   `json:"name"`
	Region       string   `json:"region"`
	Country      string   `json:"country"`
	ISP          string   `json:"isp"`
	Status       string   `json:"status"`
	Tags         []string `json:"tags"`
	Protocols    []string `json:"protocols"`
	CapacityMbps int      `json:"capacity_mbps"`
	Description  string   `json:"description"`
}

// AdminUpdateNodeRequest 管理端更新节点请求。
type AdminUpdateNodeRequest struct {
	NodeID       uint64   `path:"id"`
	Name         *string  `json:"name,optional"`
	Region       *string  `json:"region,optional"`
	Country      *string  `json:"country,optional"`
	ISP          *string  `json:"isp,optional"`
	Status       *string  `json:"status,optional"`
	Tags         []string `json:"tags,optional"`
	Protocols    []string `json:"protocols,optional"`
	CapacityMbps *int     `json:"capacity_mbps,optional"`
	Description  *string  `json:"description,optional"`
}

// AdminDisableNodeRequest 管理端禁用节点请求。
type AdminDisableNodeRequest struct {
	NodeID uint64 `path:"id"`
}

// AdminNodeListResponse 节点列表响应。
type AdminNodeListResponse struct {
	Nodes      []NodeSummary  `json:"nodes"`
	Pagination PaginationMeta `json:"pagination"`
}

// AdminNodeKernelsRequest 请求节点协议配置。
type AdminNodeKernelsRequest struct {
	NodeID uint64 `path:"id"`
}

// NodeKernelSummary 节点协议详情。
type NodeKernelSummary struct {
	Protocol     string         `json:"protocol"`
	Endpoint     string         `json:"endpoint"`
	Revision     string         `json:"revision"`
	Status       string         `json:"status"`
	Config       map[string]any `json:"config"`
	LastSyncedAt int64          `json:"last_synced_at"`
}

// AdminNodeKernelResponse 节点协议列表返回。
type AdminNodeKernelResponse struct {
	NodeID  uint64              `json:"node_id"`
	Kernels []NodeKernelSummary `json:"kernels"`
}

// AdminUpsertNodeKernelRequest 管理端配置节点内核端点。
type AdminUpsertNodeKernelRequest struct {
	NodeID       uint64         `path:"id"`
	Protocol     string         `json:"protocol"`
	Endpoint     string         `json:"endpoint"`
	Revision     string         `json:"revision"`
	Status       string         `json:"status"`
	Config       map[string]any `json:"config"`
	LastSyncedAt *int64         `json:"last_synced_at,omitempty,optional"`
}

// AdminNodeKernelUpsertResponse 节点内核端点更新响应。
type AdminNodeKernelUpsertResponse struct {
	NodeID uint64            `json:"node_id"`
	Kernel NodeKernelSummary `json:"kernel"`
}

// AdminSyncNodeKernelRequest 触发节点同步请求。
type AdminSyncNodeKernelRequest struct {
	NodeID   uint64 `path:"id"`
	Protocol string `json:"protocol"`
}

// AdminSyncNodeKernelResponse 返回最新同步信息。
type AdminSyncNodeKernelResponse struct {
	NodeID   uint64 `json:"node_id"`
	Protocol string `json:"protocol"`
	Revision string `json:"revision"`
	SyncedAt int64  `json:"synced_at"`
	Message  string `json:"message"`
}

// AdminListProtocolConfigsRequest 管理端协议配置列表请求。
type AdminListProtocolConfigsRequest struct {
	Page      int    `form:"page,optional" json:"page,optional"`
	PerPage   int    `form:"per_page,optional" json:"per_page,optional"`
	Sort      string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Query     string `form:"q,optional" json:"q,optional"`
	Protocol  string `form:"protocol,optional" json:"protocol,optional"`
	Status    string `form:"status,optional" json:"status,optional"`
}

// ProtocolConfigSummary 协议配置摘要。
type ProtocolConfigSummary struct {
	ID          uint64         `json:"id"`
	Name        string         `json:"name"`
	Protocol    string         `json:"protocol"`
	Status      string         `json:"status"`
	Tags        []string       `json:"tags"`
	Description string         `json:"description"`
	Profile     map[string]any `json:"profile"`
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
}

// AdminProtocolConfigListResponse 协议配置列表响应。
type AdminProtocolConfigListResponse struct {
	Configs    []ProtocolConfigSummary `json:"configs"`
	Pagination PaginationMeta          `json:"pagination"`
}

// AdminCreateProtocolConfigRequest 管理端创建协议配置请求。
type AdminCreateProtocolConfigRequest struct {
	Name        string         `json:"name"`
	Protocol    string         `json:"protocol"`
	Status      string         `json:"status"`
	Tags        []string       `json:"tags"`
	Description string         `json:"description"`
	Profile     map[string]any `json:"profile"`
}

// AdminUpdateProtocolConfigRequest 管理端更新协议配置请求。
type AdminUpdateProtocolConfigRequest struct {
	ConfigID    uint64         `path:"id"`
	Name        *string        `json:"name,optional"`
	Protocol    *string        `json:"protocol,optional"`
	Status      *string        `json:"status,optional"`
	Tags        []string       `json:"tags,optional"`
	Description *string        `json:"description,optional"`
	Profile     map[string]any `json:"profile,optional"`
}

// AdminDeleteProtocolConfigRequest 删除协议配置请求。
type AdminDeleteProtocolConfigRequest struct {
	ConfigID uint64 `path:"id"`
}

// AdminListProtocolBindingsRequest 管理端协议绑定列表请求。
type AdminListProtocolBindingsRequest struct {
	Page             int     `form:"page,optional" json:"page,optional"`
	PerPage          int     `form:"per_page,optional" json:"per_page,optional"`
	Sort             string  `form:"sort,optional" json:"sort,optional"`
	Direction        string  `form:"direction,optional" json:"direction,optional"`
	Query            string  `form:"q,optional" json:"q,optional"`
	Status           string  `form:"status,optional" json:"status,optional"`
	Protocol         string  `form:"protocol,optional" json:"protocol,optional"`
	NodeID           *uint64 `form:"node_id,optional" json:"node_id,optional"`
	ProtocolConfigID *uint64 `form:"protocol_config_id,optional" json:"protocol_config_id,optional"`
}

// ProtocolBindingSummary 协议绑定摘要。
type ProtocolBindingSummary struct {
	ID               uint64         `json:"id"`
	Name             string         `json:"name"`
	NodeID           uint64         `json:"node_id"`
	NodeName         string         `json:"node_name"`
	ProtocolConfigID uint64         `json:"protocol_config_id"`
	Protocol         string         `json:"protocol"`
	Role             string         `json:"role"`
	Listen           string         `json:"listen"`
	Connect          string         `json:"connect"`
	Status           string         `json:"status"`
	KernelID         string         `json:"kernel_id"`
	SyncStatus       string         `json:"sync_status"`
	HealthStatus     string         `json:"health_status"`
	LastSyncedAt     int64          `json:"last_synced_at"`
	LastHeartbeatAt  int64          `json:"last_heartbeat_at"`
	LastSyncError    string         `json:"last_sync_error"`
	Tags             []string       `json:"tags"`
	Description      string         `json:"description"`
	Metadata         map[string]any `json:"metadata"`
	CreatedAt        int64          `json:"created_at"`
	UpdatedAt        int64          `json:"updated_at"`
}

// AdminProtocolBindingListResponse 协议绑定列表响应。
type AdminProtocolBindingListResponse struct {
	Bindings   []ProtocolBindingSummary `json:"bindings"`
	Pagination PaginationMeta           `json:"pagination"`
}

// AdminCreateProtocolBindingRequest 创建协议绑定请求。
type AdminCreateProtocolBindingRequest struct {
	Name             string         `json:"name"`
	NodeID           uint64         `json:"node_id"`
	ProtocolConfigID uint64         `json:"protocol_config_id"`
	Role             string         `json:"role"`
	Listen           string         `json:"listen"`
	Connect          string         `json:"connect"`
	Status           string         `json:"status"`
	KernelID         string         `json:"kernel_id"`
	Tags             []string       `json:"tags"`
	Description      string         `json:"description"`
	Metadata         map[string]any `json:"metadata"`
}

// AdminUpdateProtocolBindingRequest 更新协议绑定请求。
type AdminUpdateProtocolBindingRequest struct {
	BindingID        uint64         `path:"id"`
	Name             *string        `json:"name,optional"`
	NodeID           *uint64        `json:"node_id,optional"`
	ProtocolConfigID *uint64        `json:"protocol_config_id,optional"`
	Role             *string        `json:"role,optional"`
	Listen           *string        `json:"listen,optional"`
	Connect          *string        `json:"connect,optional"`
	Status           *string        `json:"status,optional"`
	KernelID         *string        `json:"kernel_id,optional"`
	SyncStatus       *string        `json:"sync_status,optional"`
	HealthStatus     *string        `json:"health_status,optional"`
	LastSyncedAt     *int64         `json:"last_synced_at,omitempty,optional"`
	LastHeartbeatAt  *int64         `json:"last_heartbeat_at,omitempty,optional"`
	LastSyncError    *string        `json:"last_sync_error,optional"`
	Tags             []string       `json:"tags,optional"`
	Description      *string        `json:"description,optional"`
	Metadata         map[string]any `json:"metadata,optional"`
}

// AdminDeleteProtocolBindingRequest 删除协议绑定请求。
type AdminDeleteProtocolBindingRequest struct {
	BindingID uint64 `path:"id"`
}

// ProtocolBindingSyncResult 单条协议下发结果。
type ProtocolBindingSyncResult struct {
	BindingID uint64 `json:"binding_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	SyncedAt  int64  `json:"synced_at"`
}

// AdminSyncProtocolBindingRequest 触发单条协议下发。
type AdminSyncProtocolBindingRequest struct {
	BindingID uint64 `path:"id"`
}

// AdminSyncProtocolBindingsRequest 批量协议下发请求。
type AdminSyncProtocolBindingsRequest struct {
	BindingIDs []uint64 `json:"binding_ids"`
	NodeIDs    []uint64 `json:"node_ids"`
}

// AdminSyncProtocolBindingsResponse 批量协议下发响应。
type AdminSyncProtocolBindingsResponse struct {
	Results []ProtocolBindingSyncResult `json:"results"`
}

// AdminListSubscriptionTemplatesRequest 管理端模板列表查询。
type AdminListSubscriptionTemplatesRequest struct {
	Page          int    `form:"page,optional" json:"page,optional"`
	PerPage       int    `form:"per_page,optional" json:"per_page,optional"`
	Sort          string `form:"sort,optional" json:"sort,optional"`
	Direction     string `form:"direction,optional" json:"direction,optional"`
	Query         string `form:"q,optional" json:"q,optional"`
	ClientType    string `form:"client_type,optional" json:"client_type,optional"`
	Format        string `form:"format,optional" json:"format,optional"`
	IncludeDrafts bool   `form:"include_drafts,optional" json:"include_drafts,optional"`
}

// TemplateVariable 模板变量描述。
type TemplateVariable struct {
	ValueType    string `json:"value_type"`
	Required     bool   `json:"required"`
	Description  string `json:"description"`
	DefaultValue any    `json:"default_value"`
}

// SubscriptionTemplateSummary 模板摘要信息。
type SubscriptionTemplateSummary struct {
	ID              uint64                      `json:"id"`
	Name            string                      `json:"name"`
	Description     string                      `json:"description"`
	ClientType      string                      `json:"client_type"`
	Format          string                      `json:"format"`
	Content         string                      `json:"content,omitempty"`
	Variables       map[string]TemplateVariable `json:"variables"`
	IsDefault       bool                        `json:"is_default"`
	Version         uint32                      `json:"version"`
	UpdatedAt       int64                       `json:"updated_at"`
	PublishedAt     int64                       `json:"published_at"`
	LastPublishedBy string                      `json:"last_published_by"`
}

// AdminSubscriptionTemplateListResponse 模板列表。
type AdminSubscriptionTemplateListResponse struct {
	Templates  []SubscriptionTemplateSummary `json:"templates"`
	Pagination PaginationMeta                `json:"pagination"`
}

// AdminCreateSubscriptionTemplateRequest 创建模板。
type AdminCreateSubscriptionTemplateRequest struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	ClientType  string                      `json:"client_type"`
	Format      string                      `json:"format"`
	Content     string                      `json:"content"`
	Variables   map[string]TemplateVariable `json:"variables"`
	IsDefault   bool                        `json:"is_default"`
}

// AdminUpdateSubscriptionTemplateRequest 更新模板。
type AdminUpdateSubscriptionTemplateRequest struct {
	TemplateID  uint64                      `path:"id"`
	Name        *string                     `json:"name,optional"`
	Description *string                     `json:"description,optional"`
	Format      *string                     `json:"format,optional"`
	Content     *string                     `json:"content,optional"`
	Variables   map[string]TemplateVariable `json:"variables,optional"`
	IsDefault   *bool                       `json:"is_default,optional"`
}

// AdminPublishSubscriptionTemplateRequest 发布模板。
type AdminPublishSubscriptionTemplateRequest struct {
	TemplateID uint64 `path:"id"`
	Changelog  string `json:"changelog"`
	Operator   string `json:"operator"`
}

// AdminPublishSubscriptionTemplateResponse 发布结果。
type AdminPublishSubscriptionTemplateResponse struct {
	Template SubscriptionTemplateSummary      `json:"template"`
	History  SubscriptionTemplateHistoryEntry `json:"history"`
}

// SubscriptionTemplateHistoryEntry 模板历史条目。
type SubscriptionTemplateHistoryEntry struct {
	Version     uint32                      `json:"version"`
	Changelog   string                      `json:"changelog"`
	PublishedAt int64                       `json:"published_at"`
	PublishedBy string                      `json:"published_by"`
	Variables   map[string]TemplateVariable `json:"variables"`
}

// AdminSubscriptionTemplateHistoryRequest 查询历史。
type AdminSubscriptionTemplateHistoryRequest struct {
	TemplateID uint64 `path:"id"`
}

// AdminSubscriptionTemplateHistoryResponse 历史列表。
type AdminSubscriptionTemplateHistoryResponse struct {
	TemplateID uint64                             `json:"template_id"`
	History    []SubscriptionTemplateHistoryEntry `json:"history"`
}

// UserListSubscriptionsRequest 用户订阅列表查询。
type UserListSubscriptionsRequest struct {
	Page      int    `form:"page,optional" json:"page,optional"`
	PerPage   int    `form:"per_page,optional" json:"per_page,optional"`
	Sort      string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Query     string `form:"q,optional" json:"q,optional"`
	Status    string `form:"status,optional" json:"status,optional"`
}

// UserSubscriptionSummary 用户订阅摘要。
type UserSubscriptionSummary struct {
	ID                   uint64   `json:"id"`
	Name                 string   `json:"name"`
	PlanName             string   `json:"plan_name"`
	Status               string   `json:"status"`
	TemplateID           uint64   `json:"template_id"`
	AvailableTemplateIDs []uint64 `json:"available_template_ids"`
	ExpiresAt            int64    `json:"expires_at"`
	TrafficTotalBytes    int64    `json:"traffic_total_bytes"`
	TrafficUsedBytes     int64    `json:"traffic_used_bytes"`
	DevicesLimit         int      `json:"devices_limit"`
	LastRefreshedAt      int64    `json:"last_refreshed_at"`
}

// UserSubscriptionListResponse 用户订阅列表。
type UserSubscriptionListResponse struct {
	Subscriptions []UserSubscriptionSummary `json:"subscriptions"`
	Pagination    PaginationMeta            `json:"pagination"`
}

// UserSubscriptionPreviewRequest 用户订阅预览请求。
type UserSubscriptionPreviewRequest struct {
	SubscriptionID uint64 `path:"id"`
	TemplateID     uint64 `form:"template_id,optional" json:"template_id,optional"`
}

// UserSubscriptionPreviewResponse 用户订阅预览内容。
type UserSubscriptionPreviewResponse struct {
	SubscriptionID uint64 `json:"subscription_id"`
	TemplateID     uint64 `json:"template_id"`
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
	ETag           string `json:"etag"`
	GeneratedAt    int64  `json:"generated_at"`
}

// UserUpdateSubscriptionTemplateRequest 用户更新订阅模板。
type UserUpdateSubscriptionTemplateRequest struct {
	SubscriptionID uint64 `path:"id"`
	TemplateID     uint64 `json:"template_id"`
}

// UserUpdateSubscriptionTemplateResponse 更新结果。
type UserUpdateSubscriptionTemplateResponse struct {
	SubscriptionID uint64 `json:"subscription_id"`
	TemplateID     uint64 `json:"template_id"`
	UpdatedAt      int64  `json:"updated_at"`
}

// UserSubscriptionTrafficRequest 订阅流量明细查询请求。
type UserSubscriptionTrafficRequest struct {
	SubscriptionID    uint64  `path:"id"`
	Page              int     `form:"page,optional" json:"page,optional"`
	PerPage           int     `form:"per_page,optional" json:"per_page,optional"`
	Protocol          string  `form:"protocol,optional" json:"protocol,optional"`
	NodeID            *uint64 `form:"node_id,optional" json:"node_id,optional"`
	ProtocolBindingID *uint64 `form:"binding_id,optional" json:"binding_id,optional"`
	From              *int64  `form:"from,optional" json:"from,optional"`
	To                *int64  `form:"to,optional" json:"to,optional"`
}

// UserTrafficUsageRecord 流量记录摘要。
type UserTrafficUsageRecord struct {
	ID                uint64  `json:"id"`
	Protocol          string  `json:"protocol"`
	NodeID            uint64  `json:"node_id"`
	ProtocolBindingID uint64  `json:"binding_id"`
	BytesUp           int64   `json:"bytes_up"`
	BytesDown         int64   `json:"bytes_down"`
	RawBytes          int64   `json:"raw_bytes"`
	ChargedBytes      int64   `json:"charged_bytes"`
	Multiplier        float64 `json:"multiplier"`
	ObservedAt        int64   `json:"observed_at"`
}

// UserSubscriptionTrafficSummary 流量统计汇总。
type UserSubscriptionTrafficSummary struct {
	RawBytes     int64 `json:"raw_bytes"`
	ChargedBytes int64 `json:"charged_bytes"`
}

// UserSubscriptionTrafficResponse 订阅流量明细响应。
type UserSubscriptionTrafficResponse struct {
	Summary    UserSubscriptionTrafficSummary `json:"summary"`
	Records    []UserTrafficUsageRecord       `json:"records"`
	Pagination PaginationMeta                 `json:"pagination"`
}

// AdminListPlansRequest 管理端套餐列表请求参数。
type AdminListPlansRequest struct {
	Page      int    `form:"page,optional" json:"page,optional"`
	PerPage   int    `form:"per_page,optional" json:"per_page,optional"`
	Sort      string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Query     string `form:"q,optional" json:"q,optional"`
	Status    string `form:"status,optional" json:"status,optional"`
	Visible   *bool  `form:"visible,optional" json:"visible,optional"`
}

// AdminCreatePlanRequest 管理端创建套餐请求。
type AdminCreatePlanRequest struct {
	Name               string             `json:"name"`
	Slug               string             `json:"slug"`
	Description        string             `json:"description"`
	Tags               []string           `json:"tags"`
	Features           []string           `json:"features"`
	PriceCents         int64              `json:"price_cents"`
	Currency           string             `json:"currency"`
	DurationDays       int                `json:"duration_days"`
	TrafficLimitBytes  int64              `json:"traffic_limit_bytes"`
	TrafficMultipliers map[string]float64 `json:"traffic_multipliers"`
	DevicesLimit       int                `json:"devices_limit"`
	SortOrder          int                `json:"sort_order"`
	Status             string             `json:"status"`
	Visible            bool               `json:"visible"`
}

// AdminUpdatePlanRequest 管理端更新套餐请求。
type AdminUpdatePlanRequest struct {
	PlanID             uint64             `path:"id"`
	Name               *string            `json:"name,optional"`
	Slug               *string            `json:"slug,optional"`
	Description        *string            `json:"description,optional"`
	Tags               []string           `json:"tags,optional"`
	Features           []string           `json:"features,optional"`
	PriceCents         *int64             `json:"price_cents,optional"`
	Currency           *string            `json:"currency,optional"`
	DurationDays       *int               `json:"duration_days,optional"`
	TrafficLimitBytes  *int64             `json:"traffic_limit_bytes,optional"`
	TrafficMultipliers map[string]float64 `json:"traffic_multipliers,optional"`
	DevicesLimit       *int               `json:"devices_limit,optional"`
	SortOrder          *int               `json:"sort_order,optional"`
	Status             *string            `json:"status,optional"`
	Visible            *bool              `json:"visible,optional"`
}

// PlanSummary 套餐概览。
type PlanSummary struct {
	ID                 uint64                     `json:"id"`
	Name               string                     `json:"name"`
	Slug               string                     `json:"slug"`
	Description        string                     `json:"description"`
	Tags               []string                   `json:"tags"`
	Features           []string                   `json:"features"`
	BillingOptions     []PlanBillingOptionSummary `json:"billing_options"`
	PriceCents         int64                      `json:"price_cents"`
	Currency           string                     `json:"currency"`
	DurationDays       int                        `json:"duration_days"`
	TrafficLimitBytes  int64                      `json:"traffic_limit_bytes"`
	TrafficMultipliers map[string]float64         `json:"traffic_multipliers"`
	DevicesLimit       int                        `json:"devices_limit"`
	SortOrder          int                        `json:"sort_order"`
	Status             string                     `json:"status"`
	Visible            bool                       `json:"visible"`
	CreatedAt          int64                      `json:"created_at"`
	UpdatedAt          int64                      `json:"updated_at"`
}

// AdminPlanListResponse 管理端套餐列表响应。
type AdminPlanListResponse struct {
	Plans      []PlanSummary  `json:"plans"`
	Pagination PaginationMeta `json:"pagination"`
}

// AdminListAnnouncementsRequest 管理端公告列表参数。
type AdminListAnnouncementsRequest struct {
	Page      int    `form:"page,optional" json:"page,optional"`
	PerPage   int    `form:"per_page,optional" json:"per_page,optional"`
	Status    string `form:"status,optional" json:"status,optional"`
	Category  string `form:"category,optional" json:"category,optional"`
	Audience  string `form:"audience,optional" json:"audience,optional"`
	Query     string `form:"q,optional" json:"q,optional"`
	Sort      string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
}

// AdminCreateAnnouncementRequest 创建公告。
type AdminCreateAnnouncementRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	Audience  string `json:"audience"`
	IsPinned  bool   `json:"is_pinned"`
	Priority  int    `json:"priority"`
	CreatedBy string `json:"created_by"`
}

// AdminPublishAnnouncementRequest 发布公告。
type AdminPublishAnnouncementRequest struct {
	AnnouncementID uint64 `path:"id"`
	VisibleTo      int64  `json:"visible_to"`
	Operator       string `json:"operator"`
}

// AnnouncementSummary 公告信息。
type AnnouncementSummary struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Status      string `json:"status"`
	Audience    string `json:"audience"`
	IsPinned    bool   `json:"is_pinned"`
	Priority    int    `json:"priority"`
	VisibleFrom int64  `json:"visible_from"`
	VisibleTo   *int64 `json:"visible_to"`
	PublishedAt *int64 `json:"published_at"`
	PublishedBy string `json:"published_by"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// AdminAnnouncementListResponse 管理端公告响应。
type AdminAnnouncementListResponse struct {
	Announcements []AnnouncementSummary `json:"announcements"`
	Pagination    PaginationMeta        `json:"pagination"`
}

// SecuritySetting 第三方安全配置。
type SecuritySetting struct {
	ID                   uint64 `json:"id"`
	ThirdPartyAPIEnabled bool   `json:"third_party_api_enabled"`
	APIKey               string `json:"api_key"`
	APISecret            string `json:"api_secret"`
	EncryptionAlgorithm  string `json:"encryption_algorithm"`
	NonceTTLSeconds      int    `json:"nonce_ttl_seconds"`
	CreatedAt            int64  `json:"created_at"`
	UpdatedAt            int64  `json:"updated_at"`
}

// AdminSecuritySettingResponse 安全配置响应。
type AdminSecuritySettingResponse struct {
	Setting SecuritySetting `json:"setting"`
}

// AdminUpdateSecuritySettingRequest 更新安全配置。
type AdminUpdateSecuritySettingRequest struct {
	ThirdPartyAPIEnabled *bool   `json:"third_party_api_enabled,optional"`
	APIKey               *string `json:"api_key,optional"`
	APISecret            *string `json:"api_secret,optional"`
	EncryptionAlgorithm  *string `json:"encryption_algorithm,optional"`
	NonceTTLSeconds      *int    `json:"nonce_ttl_seconds,optional"`
}

// UserPlanListRequest 用户套餐列表参数。
type UserPlanListRequest struct {
	Query string `form:"q,optional" json:"q,optional"`
}

// UserPlanSummary 用户侧套餐信息。
type UserPlanSummary struct {
	ID                uint64                     `json:"id"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Features          []string                   `json:"features"`
	BillingOptions    []PlanBillingOptionSummary `json:"billing_options"`
	PriceCents        int64                      `json:"price_cents"`
	Currency          string                     `json:"currency"`
	DurationDays      int                        `json:"duration_days"`
	TrafficLimitBytes int64                      `json:"traffic_limit_bytes"`
	DevicesLimit      int                        `json:"devices_limit"`
	Tags              []string                   `json:"tags"`
}

// UserPlanListResponse 用户套餐列表。
type UserPlanListResponse struct {
	Plans []UserPlanSummary `json:"plans"`
}

// UserNodeStatusListRequest 用户节点状态列表请求。
type UserNodeStatusListRequest struct {
	Page     int    `form:"page,optional" json:"page,optional"`
	PerPage  int    `form:"per_page,optional" json:"per_page,optional"`
	Status   string `form:"status,optional" json:"status,optional"`
	Protocol string `form:"protocol,optional" json:"protocol,optional"`
}

// UserNodeKernelStatusSummary 用户端节点协议状态摘要（脱敏）。
type UserNodeKernelStatusSummary struct {
	Protocol     string `json:"protocol"`
	Status       string `json:"status"`
	LastSyncedAt int64  `json:"last_synced_at"`
}

// UserNodeProtocolStatusSummary 用户端协议运行状态摘要（脱敏）。
type UserNodeProtocolStatusSummary struct {
	BindingID       uint64 `json:"binding_id"`
	Protocol        string `json:"protocol"`
	Role            string `json:"role"`
	Status          string `json:"status"`
	HealthStatus    string `json:"health_status"`
	LastHeartbeatAt int64  `json:"last_heartbeat_at"`
}

// UserNodeStatusSummary 用户端节点运行状态摘要。
type UserNodeStatusSummary struct {
	ID               uint64                          `json:"id"`
	Name             string                          `json:"name"`
	Region           string                          `json:"region"`
	Country          string                          `json:"country"`
	ISP              string                          `json:"isp"`
	Status           string                          `json:"status"`
	Tags             []string                        `json:"tags"`
	Protocols        []string                        `json:"protocols"`
	CapacityMbps     int                             `json:"capacity_mbps"`
	Description      string                          `json:"description"`
	LastSyncedAt     int64                           `json:"last_synced_at"`
	UpdatedAt        int64                           `json:"updated_at"`
	KernelStatuses   []UserNodeKernelStatusSummary   `json:"kernel_statuses"`
	ProtocolStatuses []UserNodeProtocolStatusSummary `json:"protocol_statuses"`
}

// UserNodeStatusListResponse 用户端节点状态列表响应。
type UserNodeStatusListResponse struct {
	Nodes      []UserNodeStatusSummary `json:"nodes"`
	Pagination PaginationMeta          `json:"pagination"`
}

// UserProfile 用户基础资料。
type UserProfile struct {
	ID              uint64 `json:"id"`
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	Status          string `json:"status"`
	EmailVerifiedAt *int64 `json:"email_verified_at"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// CredentialSummary summarizes user credential metadata.
type CredentialSummary struct {
	Version      int    `json:"version"`
	Status       string `json:"status"`
	IssuedAt     int64  `json:"issued_at"`
	DeprecatedAt *int64 `json:"deprecated_at,omitempty"`
	RevokedAt    *int64 `json:"revoked_at,omitempty"`
	LastSeenAt   *int64 `json:"last_seen_at,omitempty"`
}

// UserProfileRequest 用户资料请求。
type UserProfileRequest struct{}

// UserProfileResponse 用户资料响应。
type UserProfileResponse struct {
	Profile UserProfile `json:"profile"`
}

// UserUpdateProfileRequest 用户资料更新请求。
type UserUpdateProfileRequest struct {
	DisplayName *string `json:"display_name,optional"`
}

// UserChangePasswordRequest 用户自助改密请求。
type UserChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// UserChangePasswordResponse 用户改密响应。
type UserChangePasswordResponse struct {
	Message string `json:"message"`
}

// UserRotateCredentialRequest triggers credential rotation.
type UserRotateCredentialRequest struct{}

// UserRotateCredentialResponse returns credential metadata.
type UserRotateCredentialResponse struct {
	Credential CredentialSummary `json:"credential"`
}

// UserEmailChangeCodeRequest 邮箱变更验证码请求。
type UserEmailChangeCodeRequest struct {
	Email string `json:"email"`
}

// UserEmailChangeCodeResponse 邮箱变更验证码响应。
type UserEmailChangeCodeResponse struct {
	Message string `json:"message"`
}

// UserChangeEmailRequest 用户邮箱变更请求。
type UserChangeEmailRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

// UserChangeEmailResponse 邮箱变更响应。
type UserChangeEmailResponse struct {
	Profile UserProfile `json:"profile"`
}

// UserAnnouncementListRequest 用户公告请求。
type UserAnnouncementListRequest struct {
	Audience string `form:"audience,optional" json:"audience,optional"`
	Limit    int    `form:"limit,optional" json:"limit,optional"`
}

// UserAnnouncementSummary 用户公告信息。
type UserAnnouncementSummary struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Audience    string `json:"audience"`
	IsPinned    bool   `json:"is_pinned"`
	Priority    int    `json:"priority"`
	VisibleFrom int64  `json:"visible_from"`
	VisibleTo   *int64 `json:"visible_to"`
	PublishedAt *int64 `json:"published_at"`
}

// UserAnnouncementListResponse 用户公告响应。
type UserAnnouncementListResponse struct {
	Announcements []UserAnnouncementSummary `json:"announcements"`
}

// UserBalanceRequest 用户余额请求。
type UserBalanceRequest struct {
	Page      int    `form:"page,optional" json:"page,optional"`
	PerPage   int    `form:"per_page,optional" json:"per_page,optional"`
	EntryType string `form:"entry_type,optional" json:"entry_type,optional"`
}

// BalanceTransactionSummary 用户余额流水。
type BalanceTransactionSummary struct {
	ID                uint64         `json:"id"`
	EntryType         string         `json:"entry_type"`
	AmountCents       int64          `json:"amount_cents"`
	Currency          string         `json:"currency"`
	BalanceAfterCents int64          `json:"balance_after_cents"`
	Reference         string         `json:"reference"`
	Description       string         `json:"description"`
	Metadata          map[string]any `json:"metadata"`
	CreatedAt         int64          `json:"created_at"`
}

// UserBalanceResponse 用户余额详情。
type UserBalanceResponse struct {
	UserID       uint64                      `json:"user_id"`
	BalanceCents int64                       `json:"balance_cents"`
	Currency     string                      `json:"currency"`
	UpdatedAt    int64                       `json:"updated_at"`
	Transactions []BalanceTransactionSummary `json:"transactions"`
	Pagination   PaginationMeta              `json:"pagination"`
}

// BalanceSnapshot 余额快照。
type BalanceSnapshot struct {
	UserID       uint64 `json:"user_id"`
	BalanceCents int64  `json:"balance_cents"`
	Currency     string `json:"currency"`
	UpdatedAt    int64  `json:"updated_at"`
}
