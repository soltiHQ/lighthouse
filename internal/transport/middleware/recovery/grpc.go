package recovery

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Unary returns unary gRPC interceptor that recovers from panics and logs them.
func Unary(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		defer func() {
			if rec := recover(); rec != nil {
				duration := time.Since(start)

				ev := logger.Error().
					Interface("panic", rec).
					Str("method", info.FullMethod).
					Str("remote_addr", remoteAddrGRPC(ctx)).
					Str("status", codes.Internal.String()).
					Int64("duration_ms", duration.Milliseconds())

				ev = withRequestID(ctx, ev)
				ev.Msg("grpc panic recovered")

				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
