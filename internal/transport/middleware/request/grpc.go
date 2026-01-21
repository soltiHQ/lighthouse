package request

import (
	"context"

	"github.com/soltiHQ/control-plane/internal/transportctx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryRequestID attaches a request ID to the context and propagates it via gRPC metadata.
func UnaryRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var requestID string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(mdRequestIDKey); len(vals) > 0 && vals[0] != "" {
				requestID = vals[0]
			}
		}
		if requestID == "" {
			requestID = newRequestID()
		}
		ctx = metadata.AppendToOutgoingContext(
			transportctx.WithRequestID(ctx, requestID),
			mdRequestIDKey,
			requestID,
		)
		return handler(ctx, req)
	}
}
