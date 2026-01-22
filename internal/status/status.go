package status

const (
	UserStatusUnknown  = 0
	UserStatusActive   = 1
	UserStatusPending  = 2
	UserStatusDisabled = 3
)

const (
	UserCredentialStatusUnknown    = 0
	UserCredentialStatusActive     = 1
	UserCredentialStatusDeprecated = 2
	UserCredentialStatusRevoked    = 3
)

const (
	AnnouncementStatusUnknown   = 0
	AnnouncementStatusDraft     = 1
	AnnouncementStatusPublished = 2
	AnnouncementStatusArchived  = 3
)

const (
	CouponStatusUnknown  = 0
	CouponStatusActive   = 1
	CouponStatusDisabled = 2
)

const (
	CouponRedemptionStatusUnknown  = 0
	CouponRedemptionStatusReserved = 1
	CouponRedemptionStatusApplied  = 2
	CouponRedemptionStatusReleased = 3
)

const (
	PlanStatusUnknown  = 0
	PlanStatusDraft    = 1
	PlanStatusActive   = 2
	PlanStatusArchived = 3
)

const (
	PlanBillingOptionStatusUnknown  = 0
	PlanBillingOptionStatusDraft    = 1
	PlanBillingOptionStatusActive   = 2
	PlanBillingOptionStatusArchived = 3
)

const (
	SubscriptionStatusUnknown  = 0
	SubscriptionStatusActive   = 1
	SubscriptionStatusDisabled = 2
	SubscriptionStatusExpired  = 3
)

const (
	NodeStatusUnknown     = 0
	NodeStatusOnline      = 1
	NodeStatusOffline     = 2
	NodeStatusMaintenance = 3
	NodeStatusDisabled    = 4
)

const (
	NodeKernelStatusUnknown    = 0
	NodeKernelStatusConfigured = 1
	NodeKernelStatusSynced     = 2
)

const (
	ProtocolBindingStatusUnknown  = 0
	ProtocolBindingStatusActive   = 1
	ProtocolBindingStatusDisabled = 2
)

const (
	ProtocolBindingSyncStatusUnknown = 0
	ProtocolBindingSyncStatusPending = 1
	ProtocolBindingSyncStatusSynced  = 2
	ProtocolBindingSyncStatusError   = 3
)

const (
	ProtocolBindingHealthStatusUnknown   = 0
	ProtocolBindingHealthStatusHealthy   = 1
	ProtocolBindingHealthStatusDegraded  = 2
	ProtocolBindingHealthStatusUnhealthy = 3
	ProtocolBindingHealthStatusOffline   = 4
)

const (
	ProtocolEntryStatusUnknown  = 0
	ProtocolEntryStatusActive   = 1
	ProtocolEntryStatusDisabled = 2
)

const (
	OrderStatusUnknown          = 0
	OrderStatusPendingPayment   = 1
	OrderStatusPaid             = 2
	OrderStatusPaymentFailed    = 3
	OrderStatusCancelled        = 4
	OrderStatusPartiallyRefunded = 5
	OrderStatusRefunded         = 6
)

const (
	OrderPaymentStatusUnknown   = 0
	OrderPaymentStatusPending   = 1
	OrderPaymentStatusSucceeded = 2
	OrderPaymentStatusFailed    = 3
)

const (
	SyncResultStatusUnknown = 0
	SyncResultStatusSynced  = 1
	SyncResultStatusError   = 2
	SyncResultStatusSkipped = 3
)

const (
	NodeSyncResultStatusUnknown = 0
	NodeSyncResultStatusOnline  = 1
	NodeSyncResultStatusOffline = 2
	NodeSyncResultStatusSkipped = 3
	NodeSyncResultStatusError   = 4
)
