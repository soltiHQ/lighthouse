package apimap

import (
	v1 "github.com/soltiHQ/control-plane/api/v1"
	"github.com/soltiHQ/control-plane/domain/model"
)

func Session(s *model.Session) v1.Session {
	if s == nil {
		return v1.Session{}
	}

	return v1.Session{
		CreatedAt: s.CreatedAt(),
		UpdatedAt: s.UpdatedAt(),
		ExpiresAt: s.ExpiresAt(),
		RevokedAt: s.RevokedAt(),

		AuthKind:     string(s.AuthKind()),
		ID:           s.ID(),
		UserID:       s.UserID(),
		CredentialID: s.CredentialID(),

		Revoked: s.Revoked(),
	}
}
