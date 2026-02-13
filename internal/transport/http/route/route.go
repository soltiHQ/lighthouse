package route

import (
	"net/http"

	"github.com/soltiHQ/control-plane/domain/kind"
)

type BaseMW func(http.Handler) http.Handler

type PermMW func(kind.Permission) BaseMW

// Chain applies middleware in order: Chain(h, a, b) => a(b(h))
func Chain(h http.Handler, mws ...BaseMW) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] == nil {
			continue
		}
		h = mws[i](h)
	}
	return h
}

// Handle registers a pattern with optional middleware.
func Handle(mux *http.ServeMux, pattern string, h http.Handler, mws ...BaseMW) {
	mux.Handle(pattern, Chain(h, mws...))
}

// HandleFunc registers a handler func with optional middleware.
func HandleFunc(mux *http.ServeMux, pattern string, fn http.HandlerFunc, mws ...BaseMW) {
	mux.Handle(pattern, Chain(fn, mws...))
}
