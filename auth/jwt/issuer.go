package jwt

import (
	"context"

	"github.com/soltiHQ/control-plane/auth"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Issuer implements auth.Issuer using HMAC-signed JWT tokens (HS256).
type Issuer struct {
	secret []byte
}

// NewIssuer creates a new JWT issuer with the provided configuration.
func NewIssuer(secret []byte) *Issuer {
	return &Issuer{secret: secret}
}

// Issue signs and returns a JWT token for the given identity.
func (i *Issuer) Issue(_ context.Context, id *auth.Identity) (string, error) {
	if id == nil {
		return "", auth.ErrInvalidToken
	}

	claims := jwtlib.MapClaims{
		"iss": id.Issuer,
		"aud": id.Audience,
		"sub": id.Subject,
		"iat": id.IssuedAt.Unix(),
		"nbf": id.NotBefore.Unix(),
		"exp": id.ExpiresAt.Unix(),
		"jti": id.TokenID,

		"uid":   id.UserID,
		"perms": id.Permissions,
	}

	t := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return t.SignedString(i.secret)
}
