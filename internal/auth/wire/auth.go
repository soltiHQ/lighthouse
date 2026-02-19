package wire

import (
	"time"

	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/auth/providers"
	passwordprovider "github.com/soltiHQ/control-plane/internal/auth/providers/password"
	"github.com/soltiHQ/control-plane/internal/auth/ratelimit"
	"github.com/soltiHQ/control-plane/internal/auth/rbac"
	session2 "github.com/soltiHQ/control-plane/internal/auth/session"
	"github.com/soltiHQ/control-plane/internal/auth/token"
	"github.com/soltiHQ/control-plane/internal/auth/token/jwt"
	"github.com/soltiHQ/control-plane/internal/storage"
)

const (
	audience = "control-plane"
	issuer   = "solti"
)

// Auth is a composition root for the authentication subsystem.
//
// It wires together:
//
//   - JWT issuer and verifier (HS256)
//   - Session service (login/refresh/revoke use cases)
//   - RBAC resolver
//   - Password auth provider
//   - Login rate limiter
//
// Auth does not implement business logic itself; it aggregates fully
// configured components ready for use by HTTP/transport layers.
type Auth struct {
	// Clock used by token issuance and verification.
	Clock token.Clock

	// Limiter tracks failed login attempts and enforces temporary blocking.
	Limiter *ratelimit.Limiter

	// Session provides login, refresh, and revoke operations.
	Session *session2.Service

	// Verifier validates incoming access tokens.
	Verifier *jwt.HSVerifier
}

// NewAuth constructs a fully wired authentication stack.
//
// Contract:
//
//   - secret must be non-empty to produce valid signed tokens.
//   - aTTL controls access token lifetime.
//   - rTTL controls refresh token lifetime.
//   - wTTL controls rate-limit block window duration.
//   - attemptLimit defines the number of failed attempts before blocking.
//   - HS256 is used for signing and verification.
//   - Refresh token rotation is enabled by default.
func NewAuth(storage storage.Storage, secret string, aTTL, rTTL, wTTL time.Duration, attemptLimit int) *Auth {
	var (
		clock   = token.RealClock()
		secretb = []byte(secret)

		verifier = jwt.NewHSVerifier(issuer, audience, secretb, clock)
		issuerHS = jwt.NewHSIssuer(secretb, clock)
		resolver = rbac.NewResolver(storage)

		sesCfg = session2.Config{
			Audience:      audience,
			Issuer:        issuer,
			AccessTTL:     aTTL,
			RefreshTTL:    rTTL,
			RotateRefresh: true,
		}
	)
	return &Auth{
		Clock:    clock,
		Verifier: verifier,
		Session: session2.New(
			storage,
			issuerHS,
			clock,
			sesCfg,
			resolver,
			map[kind.Auth]providers.Provider{
				kind.Password: passwordprovider.New(storage),
			},
		),
		Limiter: ratelimit.New(ratelimit.Config{
			MaxAttempts: attemptLimit,
			BlockWindow: wTTL,
		}),
	}
}
