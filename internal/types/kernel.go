package types

// KernelTrafficReportRequest carries kernel traffic observations.
type KernelTrafficReportRequest struct {
	Records []KernelTrafficRecord `json:"records"`
}

// KernelTrafficRecord represents a single traffic usage entry.
type KernelTrafficRecord struct {
	UserID            uint64 `json:"user_id"`
	SubscriptionID    uint64 `json:"subscription_id"`
	Protocol          string `json:"protocol"`
	NodeID            uint64 `json:"node_id"`
	ProtocolBindingID uint64 `json:"binding_id"`
	BytesUp           int64  `json:"bytes_up"`
	BytesDown         int64  `json:"bytes_down"`
	ObservedAt        int64  `json:"observed_at"`
}

// KernelTrafficIngestResponse acknowledges traffic ingestion.
type KernelTrafficIngestResponse struct {
	Accepted int `json:"accepted"`
	Failed   int `json:"failed"`
}

// KernelNodeEventRequest represents a node event notification.
type KernelNodeEventRequest struct {
	Event      string `json:"event"`
	ID         string `json:"id"`
	NodeID     string `json:"node_id"`
	Protocol   string `json:"protocol"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	ObservedAt int64  `json:"observed_at"`
}

// KernelNodeEventResponse acknowledges event handling.
type KernelNodeEventResponse struct {
	Status string `json:"status"`
}
