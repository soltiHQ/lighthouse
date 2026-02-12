package backend

import (
	"context"

	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/auth"
	"github.com/soltiHQ/control-plane/internal/auth/identity"
	"github.com/soltiHQ/control-plane/internal/auth/wire"
)

// Login implements shared login use-case.
type Login struct {
	auth *wire.Auth
}

// NewLogin creates a new Login use-case.
func NewLogin(authSvc *wire.Auth) *Login {
	if authSvc == nil {
		panic("backend: nil auth service")
	}
	return &Login{auth: authSvc}
}

// Do authenticates a user and returns issued tokens and identity.
func (x *Login) Do(ctx context.Context, subject, password, key string) (identity *identity.Identity, accessToken, refreshToken, sessionID string, err error) {
	if subject == "" || password == "" {
		return nil, "", "", "", auth.ErrInvalidRequest
	}

	now := x.auth.Clock.Now()

	if x.auth.Limiter != nil && key != "" {
		if err := x.auth.Limiter.Check(key, now); err != nil {
			return nil, "", "", "", err
		}
	}

	pair, id, err := x.auth.Session.Login(ctx, kind.Password, subject, password)
	if err != nil {
		if x.auth.Limiter != nil && key != "" {
			x.auth.Limiter.RecordFailure(key, now)
		}
		return nil, "", "", "", err
	}

	if x.auth.Limiter != nil && key != "" {
		x.auth.Limiter.Reset(key)
	}

	return id, pair.AccessToken, pair.RefreshToken, id.SessionID, nil
}
