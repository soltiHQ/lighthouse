# internal/transportctx

Transport-agnostic context values are shared across HTTP and gRPC layers.
  
The package owns two typed context keys — **Identity** and **RequestID** — and provides getters/setters for each.  
Because the keys are unexported structs, no other package can collide with them.

## What goes into context
| Value                 | Writer                             | Reader                               |
|-----------------------|------------------------------------|--------------------------------------|
| `*identity.Identity`  | Auth middleware / interceptor      | Handlers, loggers, permission checks |
| `string` (request ID) | RequestID middleware / interceptor | Loggers, error responders            |

## Request lifecycle

```text
  incoming request (HTTP or gRPC)
    │
    ▼
  ┌──────────────────────────────┐
  │  RequestID middleware        │  WithRequestID(ctx, rid)
  │  ── extract or generate ──   │
  └──────────────┬───────────────┘
                 │
                 ▼
  ┌──────────────────────────────┐
  │  Auth middleware             │  WithIdentity(ctx, id)
  │  ── verify token ──          │
  └──────────────┬───────────────┘
                 │
                 ▼
  ┌──────────────────────────────┐
  │  Handler / Interceptor       │  Identity(ctx), RequestID(ctx)
  │  ── business logic ──        │
  └──────────────────────────────┘
```

## Why a separate package
HTTP middleware lives in `internal/transport/http/middleware`, gRPC interceptors in `internal/transport/grpc/interceptors`.
Both need to write and read the same context values.  
Putting the keys here avoids a circular import:
```text
  transport/http/middleware ──┐
                              ├──→ transportctx ←── handler/api.go
  transport/grpc/interceptors ┘                 ←── handler/ui.go
                                                ←── loggers, responders
```
