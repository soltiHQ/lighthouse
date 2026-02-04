package models

import (
	"time"

	"github.com/soltiHQ/control-plane/domain"
)

// Agent represents a registered agent.
type Agent struct {
	UpdatedAt time.Time `json:"updated_at"`
	Uptime    int64     `json:"uptime"`

	ID       string `json:"id"`
	OS       string `json:"os,omitempty"`
	Arch     string `json:"arch,omitempty"`
	Endpoint string `json:"endpoint"`
	Platform string `json:"platform,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}

// NewAgent creates a new agent model.
func NewAgent(a *domain.AgentModel) *Agent {
	if a == nil {
		return nil
	}
	out := &Agent{
		UpdatedAt: a.UpdatedAt(),
		Uptime:    a.UptimeSeconds(),
		ID:        a.ID(),
		OS:        a.OS(),
		Arch:      a.Arch(),
		Endpoint:  a.Endpoint(),
		Platform:  a.Platform(),
		Metadata:  a.MetadataAll(),
		Labels:    a.LabelsAll(),
	}
	return out
}
