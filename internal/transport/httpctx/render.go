package httpctx

import "context"

// RenderMode defines how an HTTP response should be rendered.
type RenderMode int

const (
	RenderPage RenderMode = iota
	RenderBlock
)

type renderModeKey struct{}

// WithRenderMode stores response render mode in ctx.
func WithRenderMode(ctx context.Context, m RenderMode) context.Context {
	return context.WithValue(ctx, renderModeKey{}, m)
}

// Mode returns render mode from ctx, or default RenderPage.
func Mode(ctx context.Context) RenderMode {
	if m, ok := ctx.Value(renderModeKey{}).(RenderMode); ok {
		return m
	}
	return RenderPage
}
