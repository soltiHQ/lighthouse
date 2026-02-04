package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/soltiHQ/control-plane/auth"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Verifier implements auth.Verifier for HMAC-signed JWT tokens (HS256).
type Verifier struct {
	issuer   string
	audience string
	secret   []byte
}

// NewVerifier creates a new JWT verifier.
func NewVerifier(issuer, audience string, secret []byte) *Verifier {
	return &Verifier{
		issuer:   issuer,
		audience: audience,
		secret:   secret,
	}
}

// Verify parses and validates a raw JWT token string.
func (v *Verifier) Verify(_ context.Context, rawToken string) (*auth.Identity, error) {
	if rawToken == "" {
		return nil, auth.ErrInvalidToken
	}

	tok, err := jwtlib.Parse(rawToken, func(t *jwtlib.Token) (any, error) {
		if t.Method == nil || t.Method.Alg() != jwtlib.SigningMethodHS256.Alg() {
			return nil, auth.ErrInvalidToken
		}
		return v.secret, nil
	},
		jwtlib.WithAudience(v.audience),
		jwtlib.WithIssuer(v.issuer),
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwtlib.ErrTokenExpired),
			errors.Is(err, jwtlib.ErrTokenNotValidYet):
			return nil, auth.ErrExpiredToken
		default:
			return nil, auth.ErrInvalidToken
		}
	}
	if tok == nil || !tok.Valid {
		return nil, auth.ErrInvalidToken
	}

	mc, ok := tok.Claims.(jwtlib.MapClaims)
	if !ok {
		return nil, auth.ErrInvalidToken
	}

	id := &auth.Identity{
		RawToken:  rawToken,
		Issuer:    v.issuer,
		Audience:  []string{v.audience},
		Subject:   stringFromClaim(mc["sub"]),
		UserID:    stringFromClaim(mc["uid"]),
		TokenID:   stringFromClaim(mc["jti"]),
		IssuedAt:  time.Unix(int64FromClaim(mc["iat"]), 0),
		NotBefore: time.Unix(int64FromClaim(mc["nbf"]), 0),
		ExpiresAt: time.Unix(int64FromClaim(mc["exp"]), 0),
	}

	if perms, ok := mc["perms"].([]any); ok {
		id.Permissions = make([]string, 0, len(perms))
		for _, p := range perms {
			if s, ok := p.(string); ok {
				id.Permissions = append(id.Permissions, s)
			}
		}
	}
	return id, nil
}

func stringFromClaim(v any) string {
	s, _ := v.(string)
	return s
}

func int64FromClaim(v any) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int64:
		return x
	case int:
		return int64(x)
	default:
		return 0
	}
}
