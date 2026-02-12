package sessions

import "github.com/soltiHQ/control-plane/domain/model"

func toView(s *model.Session) View {
	if s == nil {
		return View{}
	}
	return View{
		CreatedAt: s.CreatedAt(),
		UpdatedAt: s.UpdatedAt(),
		ExpiresAt: s.ExpiresAt(),
		RevokedAt: s.RevokedAt(),

		AuthKind:     s.AuthKind(),
		ID:           s.ID(),
		UserID:       s.UserID(),
		CredentialID: s.CredentialID(),

		Revoked: s.Revoked(),
	}
}
