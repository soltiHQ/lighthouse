package backend

import (
	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/auth/identity"
)

type Nav struct {
	ShowAgents bool
	ShowUsers  bool
	ShowTasks  bool
	CanAddUser bool
}

func BuildNav(id *identity.Identity) Nav {
	if id == nil {
		return Nav{}
	}

	return Nav{
		ShowAgents: hasAnySlice(id.Permissions, kind.AgentsGet, kind.AgentsEdit),
		ShowUsers:  hasAnySlice(id.Permissions, kind.UsersGet, kind.UsersAdd, kind.UsersEdit, kind.UsersDelete),
		ShowTasks:  true,
		CanAddUser: hasAnySlice(id.Permissions, kind.UsersAdd),
	}
}

func hasAnySlice(perms []kind.Permission, ps ...kind.Permission) bool {
	for _, want := range ps {
		for _, have := range perms {
			if have == want {
				return true
			}
		}
	}
	return false
}
