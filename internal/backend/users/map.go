package users

import (
	"github.com/soltiHQ/control-plane/domain/model"
)

func toView(u *model.User) View {
	if u == nil {
		return View{}
	}
	return View{
		Permissions: u.PermissionsAll(),
		RoleIDs:     u.RoleIDsAll(),

		Subject: u.Subject(),
		Email:   u.Email(),
		Name:    u.Name(),
		ID:      u.ID(),

		Disabled: u.Disabled(),
	}
}
