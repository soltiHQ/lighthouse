package models

import (
	"time"

	"github.com/soltiHQ/control-plane/domain"
)

type Agent struct {
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at,omitempty"`

	Platform string `json:"platform,omitempty"`
	OS       string `json:"os,omitempty"`
	Arch     string `json:"arch,omitempty"`

	Labels map[string]string `json:"labels,omitempty"`
}

func NewAgent(a *domain.AgentModel) *Agent {
	if a == nil {
		return nil
	}
	out := &Agent{
		ID:        a.ID(),
		UpdatedAt: a.UpdatedAt(),

		Platform: a.Platform(),
		OS:       a.OS(),
		Arch:     a.Arch(),
	}
	return out
}
