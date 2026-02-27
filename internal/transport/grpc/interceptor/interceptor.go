// Package interceptor provides unary gRPC server interceptors for the control-plane server.
//
//   - UnaryRequestID         ensures every request carries a unique ID for log correlation.
//   - UnaryAuth              verifies access tokens and stores identity in context.
//   - UnaryRequirePermission guards RPCs by checking identity permissions.
//   - UnaryLogger            structured request/response logging with zerolog.
//   - UnaryRecovery          catches panics and returns codes.Internal to the client.
package interceptor
