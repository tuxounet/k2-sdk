# Spec 02 — Components

## Overview

Components are the primary organizational unit of a k2 application. Each component groups related controllers and has its own configuration, docs, UI, and access policy.

## Component Creation

```go
component := bases.NewBaseAppComponent(
    app,              // IApp — parent application
    name,             // string — component name (used in URL routing)
    order,            // int — loading/execution priority (lower = first)
    docs,             // *swag.Spec — Swagger spec (optional)
    ui,               // *embed.FS — embedded UI (optional)
    config,           // *embed.FS — embedded YAML config (optional)
    accessPolicy,     // IAccessPolicy — "public" or "authenticated"
    controllersCtors, // []AppControllerCtor — lazy controller constructors
)
```

## Lazy Initialization

- Components are **not** instantiated at app creation time
- Component constructors are stored and invoked on the first call to `app.GetComponents()`
- Requires the kernel to be set on the app (`app.SetKernel()` must have been called)
- If kernel is nil, `GetComponents()` returns nil

## Component Lookup

- `app.GetComponent(name)` returns the first component matching by name
- Returns `nil` if no component matches

## Dynamic Component Addition

- `app.AddComponent(ctor)` appends a constructor to the list
- Used by the plugin system to inject plugin components at runtime

## Controller Hierarchy

- Each component holds a list of controller constructors
- Controllers are lazily instantiated on first `GetControllers()` call
- `GetController(name)` looks up a controller by name, returns nil if not found

## Lifecycle

| Method             | Description                                           |
| ------------------ | ----------------------------------------------------- |
| `Init()`           | No-op by default (override in custom components)      |
| `Register(router)` | No-op by default; router group passed for route setup |
| `Start()`          | No-op by default                                      |
| `Stop()`           | No-op by default                                      |

## Ordering

Components are processed in the order they are declared in the constructors list. The `order` field can be used for sorting when needed.

## Access Policy

Each component declares an access policy:

- `AccessPolicyPublic` — no authentication required
- `AccessPolicyAuthenticated` — requires authenticated session

## Sample Usage

```go
func NewSampleComponent(app types.IApp) types.IAppComponent {
    return bases.NewBaseAppComponent(app, "sample", 0, nil, nil, &componentConfig,
        types.AccessPolicyPublic,
        []types.AppControllerCtor{
            NewSampleController,
        })
}
```
