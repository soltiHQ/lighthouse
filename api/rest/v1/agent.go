package restv1

// Agent is the REST representation of a connected agent.
type Agent struct {
	Metadata map[string]string `json:"metadata,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`

	UptimeSeconds int64 `json:"uptime_seconds"`

	ID           string `json:"id"`
	Name         string `json:"name"`
	Endpoint     string `json:"endpoint"`
	EndpointType string `json:"endpoint_type"`
	APIVersion   string `json:"api_version"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	Platform     string `json:"platform"`
	Status       string `json:"status"`
	LastSeenAt   string `json:"last_seen_at,omitempty"`

	HeartbeatInterval int `json:"heartbeat_interval_s,omitempty"`
}

// AgentListResponse is the paginated list of agents.
type AgentListResponse struct {
	Items      []Agent `json:"items"`
	NextCursor string  `json:"next_cursor,omitempty"`
}

// AgentPatchLabelsRequest is the request body for patching agent labels.
type AgentPatchLabelsRequest struct {
	Labels map[string]string `json:"labels"`
}
