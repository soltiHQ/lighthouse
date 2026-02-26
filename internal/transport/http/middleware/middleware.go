// Package middleware provides HTTP middleware for the control-plane server.
//
//   - RequestID           ensures every request carries a unique ID for log correlation.
//   - Auth                verifies access tokens and stores identity in context; handles silent refresh.
//   - RequirePermission   guards routes by checking identity permissions.
//   - Negotiate           chooses Responder (JSON / HTML) and RenderMode (page / block) per request.
//   - Logger              structured request/response logging with zerolog.
//   - Recovery            catches panics and renders a 503 response.
//   - CORS                configurable cross-origin resource sharing.
package middleware
