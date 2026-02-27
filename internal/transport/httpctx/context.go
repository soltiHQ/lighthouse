// Package httpctx stores HTTP-specific request context values.
package httpctx

import (
	"context"

	"github.com/soltiHQ/control-plane/internal/transport/http/responder"
)

type (
	responderKey  struct{}
	renderModeKey struct{}
)

var fallback responder.Responder = responder.NewJSON()

// WithResponder stores the negotiated responder in ctx.
func WithResponder(ctx context.Context, r responder.Responder) context.Context {
	return context.WithValue(ctx, responderKey{}, r)
}

// Responder returns the negotiated responder from ctx.
func Responder(ctx context.Context) responder.Responder {
	if r, ok := ctx.Value(responderKey{}).(responder.Responder); ok && r != nil {
		return r
	}
	return fallback
}

// WithRenderMode stores the render mode in ctx.
func WithRenderMode(ctx context.Context, m RenderMode) context.Context {
	return context.WithValue(ctx, renderModeKey{}, m)
}

// Mode returns the render mode from ctx, defaulting to RenderPage.
func Mode(ctx context.Context) RenderMode {
	if m, ok := ctx.Value(renderModeKey{}).(RenderMode); ok {
		return m
	}
	return RenderPage
}
