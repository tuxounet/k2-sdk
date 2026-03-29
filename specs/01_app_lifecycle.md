# Spec 01 — Application Lifecycle

## Overview

The k2-sdk provides a host runtime for Go applications. An application goes through a deterministic boot sequence managed by the kernel.

## Entry Points

Six entry points are available in the root package `k`:

| Function                      | Mode            | Description                                  |
| ----------------------------- | --------------- | -------------------------------------------- |
| `HostApp(app IApp)`           | Secure (TLS)    | Hosts the application with HTTPS support     |
| `HostUnsecureApp(app IApp)`   | Unsecure (HTTP) | Hosts the application without TLS            |
| `HostProvisionOnly(app IApp)` | Provision       | Runs only the provision playbook, then exits |
| `HostTeardownOnly(app IApp)`  | Teardown        | Runs stop + teardown playbooks, then exits   |

### CLI Flags

| Flag               | Description                                                 |
| ------------------ | ----------------------------------------------------------- |
| `--force-compute`  | Bypass the checksum cache, force execution of all playbooks |
| `--provision-only` | Run only the provision playbook, then exit (no HTTP server) |
| `--teardown-only`  | Run stop + teardown playbooks, then exit (no HTTP server)   |

These flags are parsed from `os.Args` by `HostApp` and `HostUnsecureApp`.

## Boot Sequence

```
HostApp/HostUnsecureApp
  └─ NewKernelRuntime(app, version, unsecure)
       ├─ Resolve working directory (.data or RUN_DIR env var)
       ├─ Initialize logging service (early init)
       ├─ Create all kernel services in order
       └─ Store services in context
  └─ Init()       — Initialize all services sequentially
  └─ Register()   — Register routes, components, providers
  └─ Start()      — Start all services (containers, schedules, etc.)
  └─ ListenAndServe() — Block on HTTP server, handle OS signals
```

## Application Creation

Applications are created via `NewBaseApp`:

```go
app := bases.NewBaseApp(
    name,               // string — app name (trimmed)
    version,            // string — app version (trimmed)
    docs,               // *swag.Spec — Swagger docs (optional, nil to skip)
    ui,                 // *embed.FS — embedded UI assets (optional)
    config,             // *embed.FS — embedded YAML configs (optional)
    componentsCtors,    // []AppComponentCtor — lazy component constructors
    externalComponents, // *embed.FS — plugin .so files (optional)
)
```

## Shutdown

- OS signals `SIGINT` and `SIGTERM` trigger graceful shutdown
- `Stop()` is called on all services in **reverse** registration order
- Each service's `Stop()` method handles its own cleanup

## Sample Usage (from samples/simple)

```go
func main() {
    runtime.HostUnsecureApp(app.NewApp())
}
```

```go
func NewApp() types.IApp {
    return bases.NewBaseApp("sample-app", "0.0.1", nil, nil, &appConfig,
        []types.AppComponentCtor{
            components.NewSampleComponent,
        }, nil)
}
```
