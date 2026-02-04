package logger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Unary logs completed gRPC unary requests.
func Unary(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)
		duration := time.Since(start)
		code := status.Code(err)

		ev := logger.Info().
			Str("method", info.FullMethod).
			Str("remote_addr", remoteAddrGRPC(ctx)).
			Str("status", code.String()).
			Int64("duration_ms", duration.Milliseconds())

		ev = withRequestID(ctx, ev)
		ev.Msg("grpc request completed")
		return resp, err
	}
}
