package kernel

// ProtocolUpsertRequest aligns with core.yaml ProtocolUpsertRequest.
type ProtocolUpsertRequest struct {
	Listen  string      `json:"listen,omitempty"`
	Connect string      `json:"connect,omitempty"`
	Users   []User      `json:"users,omitempty"`
	Profile NodeProfile `json:"profile"`
}

// ProtocolSummary is the response payload from protocol upserts.
type ProtocolSummary struct {
	ID          string   `json:"id"`
	Role        string   `json:"role"`
	Protocol    string   `json:"protocol"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Listen      string   `json:"listen"`
	Connect     string   `json:"connect"`
	UserCount   int      `json:"user_count"`
}

// NodeProfile describes protocol profile payload for kernel.
type NodeProfile struct {
	ID          string         `json:"id"`
	Role        string         `json:"role"`
	Protocol    string         `json:"protocol"`
	Tags        []string       `json:"tags,omitempty"`
	Description string         `json:"description,omitempty"`
	Profile     map[string]any `json:"profile"`
}

// User describes kernel user payload (subset).
type User struct {
	ID       string         `json:"id"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Tags     []string       `json:"tags,omitempty"`
}

// StatusResponse carries kernel status snapshot (subset).
type StatusResponse struct {
	Snapshot RuntimeStatusSnapshot `json:"snapshot"`
}

// RuntimeStatusSnapshot contains node health list (subset).
type RuntimeStatusSnapshot struct {
	Nodes []NodeStatus `json:"nodes"`
}

// NodeStatus describes protocol node health.
type NodeStatus struct {
	ID       string          `json:"id"`
	Role     string          `json:"role"`
	Protocol string          `json:"protocol"`
	Health   NodeHealthState `json:"health"`
}

// NodeHealthState represents node health status.
type NodeHealthState struct {
	Status string `json:"status"`
}

// EventRegistrationRequest registers a node event callback.
type EventRegistrationRequest struct {
	Event    string `json:"event"`
	Callback string `json:"callback"`
	Secret   string `json:"secret,omitempty"`
}

// EventRegistrationRecord represents an event subscription record.
type EventRegistrationRecord struct {
	ID          string `json:"id"`
	Event       string `json:"event"`
	Callback    string `json:"callback"`
	CreatedAtMS int64  `json:"created_at_ms"`
}

// ServiceEventRegistrationRequest registers a service event callback.
type ServiceEventRegistrationRequest struct {
	Event    string `json:"event"`
	Callback string `json:"callback"`
	Secret   string `json:"secret,omitempty"`
}

// EventSubscriptionRecord represents a service event subscription record.
type EventSubscriptionRecord struct {
	ID          string `json:"id"`
	Event       string `json:"event"`
	Callback    string `json:"callback"`
	CreatedAtMS int64  `json:"created_at_ms"`
}
