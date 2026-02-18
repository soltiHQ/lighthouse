package policy

import "github.com/soltiHQ/control-plane/internal/auth/identity"

// UserDetail is a UI-oriented policy for the user detail page.
type UserDetail struct {
	CanEdit   bool // edit button, enable/disable, password
	CanDelete bool // delete button
}

// BuildUserDetail derives UI action flags from the authenticated identity.
func BuildUserDetail(id *identity.Identity) UserDetail {
	if id == nil {
		return UserDetail{}
	}

	perms := permSet(id)
	return UserDetail{
		CanEdit:   hasAny(perms, usersEdit),
		CanDelete: hasAny(perms, usersDelete),
	}
}
