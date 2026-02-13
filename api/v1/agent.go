package v1

type Agent struct {
	UptimeSeconds int64             `json:"uptime_seconds"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`

	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Platform string `json:"platform"`
	ID       string `json:"id"`
}

type AgentListResponse struct {
	Items      []Agent `json:"items"`
	NextCursor string  `json:"next_cursor,omitempty"`
}

type AgentPatchLabelsRequest struct {
	Labels map[string]string `json:"labels"`
}
