package policy

import (
	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/auth/identity"
)

// AgentDetail is a UI-oriented policy for the agent detail page.
type AgentDetail struct {
	CanEditLabels bool // edit labels button
}

// BuildAgentDetail derives UI action flags from the authenticated identity.
func BuildAgentDetail(id *identity.Identity) AgentDetail {
	if id == nil {
		return AgentDetail{}
	}

	perms := permSet(id)
	return AgentDetail{
		CanEditLabels: hasAny(perms, kind.AgentsEdit),
	}
}
