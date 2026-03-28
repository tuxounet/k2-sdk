# Spec 06 — Kernel Services

## Overview

The kernel provides a service registry pattern where services are stored in a Go `context.Context` and accessed by key. All services follow a common lifecycle contract.

## Service Registry

Services are stored in the kernel's root context using typed keys:

```go
type KernelServiceContextKey string

kernel.SetService(service)                              // stores by service.GetName()
kernel.GetService(KernelServiceContextKey("config"))    // retrieves from context
```

## Service Registration Order

Services are created and initialized in this fixed order:

1. `monitoring.logging` (early init — before all others)
2. `storage.paths`
3. `profile`
4. `storage.volumes`
5. `storage.stores`
6. `config`
7. `secrets`
8. `compute`
9. `plugins`
10. `app`
11. `scheduler`
12. `network.ingress.http`

**Stop order** is the reverse of this list.

## IKernelService Interface

```go
type IKernelService interface {
    ILoggable
    GetName() string
    GetKernel() IKernel
    GetConfig(key string) string
    SetConfig(key string, value string)
    GetData(key string) interface{}
    SetData(key string, value interface{})
    Init() error
    Register() error
    Start() error
    Stop() error
}
```

## BaseKernelService

Default implementation for all kernel services:

```go
base := bases.NewBaseKernelService(kernel, "service-name")
```

- Creates a sub-logger named after the service
- Initializes empty `config` (map[string]string) and `data` (map[string]interface{}) maps
- All lifecycle methods are no-ops (return nil)
- Config: `GetConfig("")` returns `""` for missing keys
- Data: `GetData("")` returns `nil` for missing keys

## Lifecycle

| Phase | Method       | Description                                |
| ----- | ------------ | ------------------------------------------ |
| 1     | `Init()`     | Initialize service state                   |
| 2     | `Register()` | Register with other services, setup routes |
| 3     | `Start()`    | Start active processing                    |
| 4     | `Stop()`     | Graceful shutdown (reverse order)          |

## Service Access Pattern

```go
kernel := controller.GetComponent().GetApp().GetKernel()
configSvc := kernel.GetService("config")
storeSvc := kernel.GetService("storage.stores")
```
