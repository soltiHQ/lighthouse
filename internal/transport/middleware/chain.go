package middleware

import (
	"net/http"

	"github.com/soltiHQ/control-plane/internal/transport/middleware/auth"
	"github.com/soltiHQ/control-plane/internal/transport/middleware/cors"
	"github.com/soltiHQ/control-plane/internal/transport/middleware/logger"
	"github.com/soltiHQ/control-plane/internal/transport/middleware/recovery"
	"github.com/soltiHQ/control-plane/internal/transport/middleware/request"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// HttpChain builds an HTTP middleware chain around the given handler according to the provided config.
func HttpChain(base http.Handler, log zerolog.Logger, cfg HttpChainConfig) http.Handler {
	h := base

	if cfg.Auth == nil || !cfg.Auth.Enabled {
		h = auth.MockHTTPIdentity()(h)
	} else {
		if cfg.Auth.Verifier == nil {
			panic("middleware: auth enabled but verifier is nil")
		}
		h = auth.HTTP(cfg.Auth.Verifier, log)(h)
	}

	if cfg.CORS != nil {
		h = cors.CORS(*cfg.CORS)(h)
	}
	if cfg.Logging {
		h = logger.HTTP(log)(h)
	}
	if cfg.RequestID {
		h = request.RequestID()(h)
	}
	if cfg.Recovery {
		h = recovery.HTTP(log)(h)
	}
	return h
}

// GrpcUnaryOptions builds gRPC server options with a chained unary interceptor according to the provided config.
func GrpcUnaryOptions(log zerolog.Logger, cfg GrpcChainConfig) []grpc.ServerOption {
	var interceptors []grpc.UnaryServerInterceptor

	if cfg.Auth == nil || !cfg.Auth.Enabled {
		interceptors = append(interceptors, auth.MockUnaryIdentity())
	} else {
		if cfg.Auth.Verifier == nil {
			panic("middleware: auth enabled but verifier is nil")
		}
		interceptors = append(interceptors, auth.Unary(cfg.Auth.Verifier, log))
	}

	if cfg.Logging {
		interceptors = append(interceptors, logger.Unary(log))
	}
	if cfg.RequestID {
		interceptors = append(interceptors, request.UnaryRequestID())
	}
	if cfg.Recovery {
		interceptors = append(interceptors, recovery.Unary(log))
	}
	return []grpc.ServerOption{grpc.ChainUnaryInterceptor(interceptors...)}
}
