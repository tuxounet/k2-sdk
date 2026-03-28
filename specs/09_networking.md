# Spec 09 — Networking (Ingress)

## Overview

The ingress service exposes the application over HTTP using the Gin web framework. It handles routing, CORS, authentication middleware, UI serving, and API documentation.

## Service Key

`"network.ingress.http"` — accessible via `kernel.GetService("network.ingress.http")`

## Route Structure

```
/{baseUrl}/
├── ui/                          # App-level UI (if provided)
├── docs/                        # Swagger documentation (if provided)
├── {component}/
│   ├── ui/                      # Component UI (if provided)
│   ├── docs/                    # Component docs (if provided)
│   └── {controller}/            # Controller routes
│       └── (user-defined routes)
```

The `{baseUrl}` is configured via `host.ingress.root`.

## CORS

CORS is enabled globally with permissive defaults via `gin-contrib/cors`.

## Authentication Middleware

Each route is associated with an access policy:

- `AccessPolicyPublic` — no authentication check
- `AccessPolicyAuthenticated` — authentication required

The auth middleware maps routes to policies and enforces them.

## Ingress Registration

External services (e.g., compute providers) can register additional ingresses:

```go
type IngressDefinition struct {
    IngressPath   string
    AccessPolicy  IAccessPolicy
    ServicePort   int
    ServiceHost   string
    RewritePath   *string
    CustomHandler func(def IngressDefinition) gin.HandlerFunc
}
```

## Configuration

| Key                  | Description                     |
| -------------------- | ------------------------------- |
| `host.ingress.root`  | Base URL path                   |
| `host.ingress.port`  | Listen port                     |
| `host.ingress.tls.*` | TLS configuration (secure mode) |

## Lifecycle

1. **Register()** — Create Gin engine, setup CORS, routes, middleware, register components
2. **Start()** — (handled by ListenAndServe)
3. **Listen()** — Start the HTTP server
4. **Stop()** — Shutdown the HTTP server
