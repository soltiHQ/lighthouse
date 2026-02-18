package policy

import "github.com/soltiHQ/control-plane/internal/auth/identity"

// UserDetail is a UI-oriented policy for the user detail page.
type UserDetail struct {
	IsSelf       bool // viewing own profile
	CanEdit      bool // edit basic fields, enable/disable, password
	CanEditRoles bool // edit roles & permissions (false for self)
	CanDelete    bool // delete button (false for self)
}

// BuildUserDetail derives UI action flags from the authenticated identity.
// targetUserID is the ID of the user being viewed.
func BuildUserDetail(id *identity.Identity, targetUserID string) UserDetail {
	if id == nil {
		return UserDetail{}
	}

	perms := permSet(id)
	self := id.UserID == targetUserID
	return UserDetail{
		IsSelf:       self,
		CanEdit:      hasAny(perms, usersEdit),
		CanEditRoles: hasAny(perms, usersEdit) && !self,
		CanDelete:    hasAny(perms, usersDelete) && !self,
	}
}
