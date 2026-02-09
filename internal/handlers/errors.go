package handlers

import (
	"net/http"
	"strings"

	"github.com/soltiHQ/control-plane/internal/transport/http/response"
	"github.com/soltiHQ/control-plane/ui/pages"
)

// Errors provides format-aware error responses (JSON or HTML)
// and acts as a catch-all wrapper for unmatched routes.
type Errors struct {
	json *response.JSONResponder
	html *response.HTMLResponder
}

// NewErrors creates a new error handler.
func NewErrors(json *response.JSONResponder, html *response.HTMLResponder) *Errors {
	return &Errors{json: json, html: html}
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
	if e.wantsJSON(r) {
		e.json.Error(w, r, http.StatusNotFound, "not found")
		return
	}
	e.html.Respond(w, r, http.StatusNotFound, &response.View{
		Component: pages.ErrorPage(
			http.StatusNotFound,
			"Page not found",
			"The page you are looking for doesn't exist or has been moved.",
		),
	})
}

// Unauthorized renders a 401 response.
// For HTML clients, HTMLResponder handles redirect to login.
func (e *Errors) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if e.wantsJSON(r) {
		e.json.Error(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	e.html.Error(w, r, http.StatusUnauthorized, "unauthorized")
}

// ServiceUnavailable renders a 503 response.
func (e *Errors) ServiceUnavailable(w http.ResponseWriter, r *http.Request) {
	if e.wantsJSON(r) {
		e.json.Error(w, r, http.StatusServiceUnavailable, "service unavailable")
		return
	}
	e.html.Respond(w, r, http.StatusServiceUnavailable, &response.View{
		Component: pages.ErrorPage(
			http.StatusServiceUnavailable,
			"Service unavailable",
			"The server is temporarily unable to handle the request. Please try again later.",
		),
	})
}

// wantsJSON returns true if the request expects a JSON response.
func (e *Errors) wantsJSON(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		return true
	}
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}
