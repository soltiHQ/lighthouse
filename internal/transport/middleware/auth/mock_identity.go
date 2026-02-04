package auth

import (
	"context"
	"net/http"
	"time"

	authcore "github.com/soltiHQ/control-plane/auth"
	"github.com/soltiHQ/control-plane/internal/transportctx"

	"google.golang.org/grpc"
)

func mockIdentity() *authcore.Identity {
	return &authcore.Identity{
		Subject:   "system:mock",
		UserID:    "system",
		Issuer:    "mock",
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Permissions: []string{
			"*",
		},
	}
}

// MockHTTPIdentity injects a system identity with full permissions.
func MockHTTPIdentity() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(transportctx.WithIdentity(r.Context(), mockIdentity())))
		})
	}
}

// MockUnaryIdentity injects a system identity for gRPC when auth is disabled.
func MockUnaryIdentity() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return handler(transportctx.WithIdentity(ctx, mockIdentity()), req)
	}
}
