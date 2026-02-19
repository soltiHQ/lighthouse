package policy

import (
	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/auth/identity"
)

// Nav is a UI-oriented navigation model derived from the current identity.
type Nav struct {
	ShowAgents bool
	ShowUsers  bool
	ShowTasks  bool
	CanAddUser bool
}

// BuildNav derives UI navigation flags from the authenticated identity.
func BuildNav(id *identity.Identity) Nav {
	if id == nil {
		return Nav{}
	}

	perms := permSet(id)
	return Nav{
		ShowAgents: hasAny(perms, kind.AgentsGet, kind.AgentsEdit),
		ShowUsers:  hasAny(perms, kind.UsersGet, kind.UsersAdd, kind.UsersEdit, kind.UsersDelete),
		ShowTasks:  true,
		CanAddUser: hasAny(perms, kind.UsersAdd),
	}
}
