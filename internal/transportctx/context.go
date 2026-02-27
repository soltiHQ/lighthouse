// Package transportctx provides transport-agnostic context values shared by
// HTTP middleware, gRPC interceptors, handlers, and loggers.
package transportctx

import (
	"context"

	"github.com/soltiHQ/control-plane/internal/auth/identity"
)

type (
	identityKey  struct{}
	requestIDKey struct{}
)

const unknownRequestID = "unknown"

// WithIdentity stores authenticated identity in ctx.
func WithIdentity(ctx context.Context, id *identity.Identity) context.Context {
	return context.WithValue(ctx, identityKey{}, id)
}

// WithRequestID stores request id in ctx.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// Identity returns identity from ctx (if any).
func Identity(ctx context.Context) (*identity.Identity, bool) {
	id, ok := ctx.Value(identityKey{}).(*identity.Identity)
	return id, ok && id != nil
}

// RequestID returns request id from ctx (if any).
func RequestID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(requestIDKey{}).(string)
	return rid, ok && rid != ""
}

// TryRequestID returns request id from ctx (if any).
func TryRequestID(ctx context.Context) string {
	if rid, ok := RequestID(ctx); ok {
		return rid
	}
	return unknownRequestID
}
