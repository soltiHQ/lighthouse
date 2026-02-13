package response

import (
	"net/http"

	"github.com/soltiHQ/control-plane/internal/transport/httpctx"
)

// ModeFromRequest derives page vs block semantics.
// Policy: HX-Request => block, else page.
func ModeFromRequest(r *http.Request) httpctx.RenderMode {
	if r != nil && r.Header.Get("HX-Request") == "true" {
		return httpctx.RenderBlock
	}
	return httpctx.RenderPage
}
