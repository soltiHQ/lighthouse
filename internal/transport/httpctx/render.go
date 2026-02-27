package httpctx

import "net/http"

// RenderMode tells response helpers whether to render a full HTML page or an HTMX fragment (block).
type RenderMode int

const (
	// RenderPage is a full page render (browser navigation).
	RenderPage RenderMode = iota
	// RenderBlock is a partial fragment render (HTMX swap).
	RenderBlock
)

// ModeFromRequest derives RenderMode from the incoming request.
// HTMX requests (HX-Request: true) produce RenderBlock, everything else — RenderPage.
func ModeFromRequest(r *http.Request) RenderMode {
	if r != nil && r.Header.Get("HX-Request") == "true" {
		return RenderBlock
	}
	return RenderPage
}
