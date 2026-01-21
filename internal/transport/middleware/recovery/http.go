package recovery

import (
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog"
)

// HTTP returns middleware that recovers from panics and logs them.
func HTTP(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				status  = http.StatusOK
				start   = time.Now()
				written int64
			)
			wrapped := httpsnoop.Wrap(w, httpsnoop.Hooks{
				WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						status = code
						next(code)
					}
				},
				Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(b []byte) (int, error) {
						n, err := next(b)
						written += int64(n)
						return n, err
					}
				},
			})

			defer func() {
				if rec := recover(); rec != nil {
					duration := time.Since(start)

					ev := logger.Error().
						Interface("panic", rec).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("remote_addr", remoteAddrHTTP(r)).
						Int("status", http.StatusInternalServerError).
						Int64("bytes", written).
						Int64("duration_ms", duration.Milliseconds())

					ev = withRequestID(r.Context(), ev)
					ev.Msg("http panic recovered")

					if status == http.StatusOK && written == 0 {
						http.Error(wrapped, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(wrapped, r)
		})
	}
}
