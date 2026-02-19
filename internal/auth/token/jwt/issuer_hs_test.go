package jwt

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/soltiHQ/control-plane/internal/auth"
	"github.com/soltiHQ/control-plane/internal/auth/identity"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time { return c.t }

func validIdentity(now time.Time) *identity.Identity {
	return &identity.Identity{
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: now.Add(time.Hour),

		Issuer:  "issuer",
		Subject: "subject",
		UserID:  "user-1",

		Audience:  []string{"aud"},
		TokenID:   "tid",
		SessionID: "sid",
	}
}

func parseNoTimeValidation(t *testing.T, tokenStr string, key []byte) jwtlib.MapClaims {
	t.Helper()

	parser := jwtlib.NewParser(jwtlib.WithoutClaimsValidation())

	parsed, err := parser.Parse(tokenStr, func(tok *jwtlib.Token) (any, error) {
		if tok.Method != jwtlib.SigningMethodHS256 {
			t.Fatalf("unexpected signing method: %v", tok.Method)
		}
		return key, nil
	})
	if err != nil {
		t.Fatalf("failed to parse/verify token: %v", err)
	}
	if !parsed.Valid {
		t.Fatalf("token not valid")
	}

	claims, ok := parsed.Claims.(jwtlib.MapClaims)
	if !ok {
		t.Fatalf("unexpected claims type: %T", parsed.Claims)
	}
	return claims
}

func TestHSIssuer_Issue_InvalidInput(t *testing.T) {
	clk := &fakeClock{t: time.Unix(100, 0)}
	iss := NewHSIssuer([]byte("secret"), clk)

	t.Run("nil identity", func(t *testing.T) {
		_, err := iss.Issue(context.Background(), nil)
		if !errors.Is(err, auth.ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		id := &identity.Identity{}
		_, err := iss.Issue(context.Background(), id)
		if !errors.Is(err, auth.ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("empty audience", func(t *testing.T) {
		id := validIdentity(clk.Now())
		id.Audience = nil
		_, err := iss.Issue(context.Background(), id)
		if !errors.Is(err, auth.ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("empty secret", func(t *testing.T) {
		iss2 := NewHSIssuer(nil, clk)
		id := validIdentity(clk.Now())
		_, err := iss2.Issue(context.Background(), id)
		if !errors.Is(err, auth.ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})
}

func TestHSIssuer_Issue_SuccessAndVerify(t *testing.T) {
	now := time.Unix(200, 0)
	clk := &fakeClock{t: now}
	secret := []byte("super-secret")

	iss := NewHSIssuer(secret, clk)

	id := validIdentity(now)
	tokenStr, err := iss.Issue(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token")
	}

	claims := parseNoTimeValidation(t, tokenStr, secret)

	if claims["iss"] != id.Issuer {
		t.Fatalf("issuer mismatch: %v", claims["iss"])
	}
	if claims["sub"] != id.Subject {
		t.Fatalf("subject mismatch: %v", claims["sub"])
	}
	if claims["jti"] != id.TokenID {
		t.Fatalf("token id mismatch: %v", claims["jti"])
	}

	aud, ok := claims["aud"]
	if !ok || aud == nil {
		t.Fatalf("missing aud claim: %#v", claims["aud"])
	}

	if claims["uid"] != id.UserID {
		t.Fatalf("user id mismatch: %v", claims["uid"])
	}
	if claims["sid"] != id.SessionID {
		t.Fatalf("session id mismatch: %v", claims["sid"])
	}
}

func TestHSIssuer_SecretIsCopied(t *testing.T) {
	now := time.Unix(300, 0)
	clk := &fakeClock{t: now}

	secret := []byte("orig-secret")
	iss := NewHSIssuer(secret, clk)

	for i := range secret {
		secret[i] = 'x'
	}

	id := validIdentity(now)
	tokenStr, err := iss.Issue(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = parseNoTimeValidation(t, tokenStr, []byte("orig-secret"))
}
