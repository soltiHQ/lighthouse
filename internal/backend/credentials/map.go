package credentials

import "github.com/soltiHQ/control-plane/domain/model"

func toView(c *model.Credential) View {
	if c == nil {
		return View{}
	}
	return View{
		Auth:   c.AuthKind(),
		UserID: c.UserID(),
		ID:     c.ID(),
	}
}
