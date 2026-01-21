package logger

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog"
)

// HTTP logs completed HTTP requests.
func HTTP(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			metrics := httpsnoop.CaptureMetrics(next, w, r)

			ev := logger.Debug().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", remoteAddrHTTP(r)).
				Int("status", metrics.Code).
				Int64("bytes", metrics.Written).
				Int64("duration_ms", metrics.Duration.Milliseconds())

			ev = withRequestID(r.Context(), ev)
			ev.Msg("http request completed")
		})
	}
}
