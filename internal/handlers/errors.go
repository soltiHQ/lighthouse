package handlers

import (
	"net/http"

	"github.com/soltiHQ/control-plane/internal/transport/http/response"
	"github.com/soltiHQ/control-plane/ui/pages"
)

// Errors provides format-aware error responses.
// Uses Responder from context (set by Negotiate middleware).
type Errors struct{}

// NewErrors creates a new error handler.
func NewErrors() *Errors {
	return &Errors{}
}

// Wrap returns a handler that delegates to mux and renders 404 for unmatched routes.
func (e *Errors) Wrap(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern == "" {
			e.NotFound(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

// NotFound renders a 404 response.
func (e *Errors) NotFound(w http.ResponseWriter, r *http.Request) {
	resp := response.FromContext(r.Context())
	resp.Respond(w, r, http.StatusNotFound, &response.View{
		Data: response.ErrorBody{Code: http.StatusNotFound, Message: "not found"},
		Component: pages.ErrorPage(
			http.StatusNotFound,
			"Page not found",
			"The page you are looking for doesn't exist or has been moved.",
		),
	})
}

// Unauthorized renders a 401 response.
func (e *Errors) Unauthorized(w http.ResponseWriter, r *http.Request) {
	resp := response.FromContext(r.Context())
	resp.Error(w, r, http.StatusUnauthorized, "unauthorized")
}

// ServiceUnavailable renders a 503 response.
func (e *Errors) ServiceUnavailable(w http.ResponseWriter, r *http.Request) {
	resp := response.FromContext(r.Context())
	resp.Respond(w, r, http.StatusServiceUnavailable, &response.View{
		Data: response.ErrorBody{Code: http.StatusServiceUnavailable, Message: "service unavailable"},
		Component: pages.ErrorPage(
			http.StatusServiceUnavailable,
			"Service unavailable",
			"The server is temporarily unable to handle the request. Please try again later.",
		),
	})
}
