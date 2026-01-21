package cors

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// CORSConfig defines CORS behavior for the HTTP server.
type CORSConfig struct {
	// AllowedOrigins is a list of allowed origins.
	// Use "*" as a single entry to allow any origin.
	AllowedOrigins []string
	// AllowedMethods is a list of HTTP methods allowed for CORS requests.
	// If empty, a default set will be used.
	AllowedMethods []string
	// AllowedHeaders is a list of request headers allowed in CORS requests.
	// If empty, the middleware may reflect Access-Control-Request-Headers.
	AllowedHeaders []string
	// ExposedHeaders is a list of response headers exposed to the browser.
	ExposedHeaders []string
	// MaxAge is how long the results of a preflight request can be cached.
	// Value is expressed as a time.Duration but will be sent as seconds.
	MaxAge time.Duration
	// AllowCredentials indicates whether credentials (cookies, auth headers)
	// are allowed in cross-site requests.
	AllowCredentials bool
}

// CORS returns a middleware that adds CORS headers to responses.
func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	allowedAllOrigins := len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if allowedAllOrigins {
				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Add("Vary", "Origin")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
			} else if slices.Contains(cfg.AllowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			} else {
				next.ServeHTTP(w, r)
				return
			}
			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				if len(cfg.AllowedMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
				} else {
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				}

				if len(cfg.AllowedHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
				} else {
					reqHeaders := r.Header.Get("Access-Control-Request-Headers")
					if reqHeaders != "" {
						w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
						w.Header().Add("Vary", "Access-Control-Request-Headers")
					}
				}

				if cfg.MaxAge > 0 {
					seconds := cfg.MaxAge / time.Second
					w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(seconds), 10))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if len(cfg.ExposedHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposedHeaders, ", "))
			}
			next.ServeHTTP(w, r)
		})
	}
}
