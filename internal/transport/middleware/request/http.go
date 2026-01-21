package request

import (
	"net/http"

	"github.com/soltiHQ/control-plane/internal/transportctx"
)

// RequestID attaches a request ID to the context and HTTP response headers.
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(headerRequestID)
			if requestID == "" {
				requestID = newRequestID()
			}

			ctx := transportctx.WithRequestID(r.Context(), requestID)
			w.Header().Set(headerRequestID, requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
