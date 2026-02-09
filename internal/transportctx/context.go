package transportctx

import (
	"context"

	"github.com/soltiHQ/control-plane/internal/auth/identity"
)

// Typed keys (unexported) prevent collisions with other context users.
type (
	identityKey  struct{}
	requestIDKey struct{}
)

// WithIdentity stores authenticated identity in ctx.
// Passing nil clears the identity value (returns a derived context anyway).
func WithIdentity(ctx context.Context, id *identity.Identity) context.Context {
	return context.WithValue(ctx, identityKey{}, id)
}

// WithRequestID stores request id in ctx.
// RequestID should be stable per request, suitable for log correlation.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// Identity returns identity from ctx (if any).
func Identity(ctx context.Context) (*identity.Identity, bool) {
	v := ctx.Value(identityKey{})
	if v == nil {
		return nil, false
	}
	id, ok := v.(*identity.Identity)
	return id, ok && id != nil
}

// RequestID returns request id from ctx (if any).
func RequestID(ctx context.Context) (string, bool) {
	v := ctx.Value(requestIDKey{})
	s, ok := v.(string)
	return s, ok && s != ""
}
