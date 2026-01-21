package request

import "github.com/segmentio/ksuid"

// headerRequestID is the HTTP header name used to propagate request IDs.
const headerRequestID = "X-Request-ID"

// mdRequestIDKey is the gRPC metadata key used to propagate request IDs.
const mdRequestIDKey = "x-request-id"

func newRequestID() string {
	return ksuid.New().String()
}
