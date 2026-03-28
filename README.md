# k2-sdk

Go SDK for building kernel-driven applications with a structured lifecycle, modular components, and pluggable infrastructure services.

## Installation

```sh
go get -v github.com/tuxounet/k2-sdk
```

Requires **Go 1.24+**.

## Overview

k2-sdk provides a host runtime that manages the full lifecycle of a Go application through a layered architecture:

```
Application
├── Components (organizational units)
│   └── Controllers (HTTP routes & scheduled tasks)
└── Kernel
    ├── config          — hierarchical YAML configuration
    ├── monitoring      — structured logging
    ├── compute         — container/k8s/playbook provisioning
    ├── storage         — object stores & volumes
    ├── network.ingress — HTTP server (Gin), CORS, routing
    ├── secrets         — encrypted credentials & SSH keys
    ├── profile         — persistent app metadata
    ├── plugins         — dynamic .so loading
    └── scheduler       — cron-based task scheduling
```

## Entry Points

| Function                 | Mode            | Description                              |
| ------------------------ | --------------- | ---------------------------------------- |
| `k.HostApp(app)`         | Secure (TLS)    | Hosts the application with HTTPS support |
| `k.HostUnsecureApp(app)` | Unsecure (HTTP) | Hosts the application without TLS        |

## Boot Sequence

```
HostApp / HostUnsecureApp
  └─ NewKernelRuntime(app, version, unsecure)
       ├─ Resolve working directory (.data or RUN_DIR env var)
       ├─ Initialize logging (early init)
       ├─ Create all kernel services in order
       └─ Store services in context
  └─ Init()           — Initialize all services sequentially
  └─ Register()       — Register routes, components, providers
  └─ Start()          — Start services (containers, schedules…)
  └─ ListenAndServe() — Block on HTTP server, handle OS signals (SIGINT/SIGTERM)
```

Shutdown calls `Stop()` on all services in **reverse** registration order.

## Quick Start

### 1. Application definition

```go
package app

import (
    "embed"

    "github.com/tuxounet/k2-sdk/bases"
    "github.com/tuxounet/k2-sdk/types"
    "myapp/components/sample"
)

//go:embed config/*.yaml
var conf embed.FS

func NewApp() types.IApp {
    return bases.NewBaseApp(
        "my-app",                          // name
        "0.1.0",                           // version
        nil,                               // *swag.Spec (optional)
        nil,                               // *embed.FS for UI (optional)
        &conf,                             // *embed.FS for YAML config (optional)
        []types.AppComponentCtor{          // component constructors
            sample.NewComponent,
        },
        nil,                               // *embed.FS for plugins (optional)
    )
}
```

### 2. Entrypoint

```go
package main

import (
    runtime "github.com/tuxounet/k2-sdk"
    "myapp/app"
)

func main() {
    runtime.HostUnsecureApp(app.NewApp())
}
```

### 3. Component

```go
package sample

import (
    "github.com/tuxounet/k2-sdk/bases"
    "github.com/tuxounet/k2-sdk/types"
)

func NewComponent(app types.IApp) types.IAppComponent {
    return bases.NewBaseAppComponent(app, "sample", 0, nil, nil, nil,
        types.AccessPolicyPublic,
        []types.AppControllerCtor{
            NewSampleController,
        })
}
```

### 4. Controller

```go
package sample

import (
    "github.com/gin-gonic/gin"
    "github.com/tuxounet/k2-sdk/bases"
    "github.com/tuxounet/k2-sdk/types"
)

type SampleController struct {
    bases.BaseAppController
}

func NewSampleController(component types.IAppComponent) types.IAppController {
    return &SampleController{
        BaseAppController: bases.NewBaseAppController(
            component, "hello", 0, nil, types.AccessPolicyPublic,
        ),
    }
}

func (c *SampleController) Register(router *gin.RouterGroup) error {
    router.GET("/sayHello", func(ctx *gin.Context) {
        ctx.JSON(200, gin.H{"message": "Hello, world!"})
    })
    return nil
}
```

Routes are mounted at `/{baseUrl}/{component}/{controller}/`.

## Kernel Services

All services follow a common lifecycle: `Init()` → `Register()` → `Start()` → `Stop()`.

Services are registered in this fixed order:

| #   | Key                    | Description                              |
| --- | ---------------------- | ---------------------------------------- |
| 1   | `monitoring.logging`   | Structured logging (early init)          |
| 2   | `storage.paths`        | Path manipulation utilities              |
| 3   | `profile`              | Persistent app metadata & properties     |
| 4   | `storage.volumes`      | Volume management                        |
| 5   | `storage.stores`       | Object storage (local / rclone backends) |
| 6   | `config`               | Hierarchical YAML config + env vars      |
| 7   | `secrets`              | Encrypted secrets & SSH key generation   |
| 8   | `compute`              | Container / Kubernetes / Playbook runner |
| 9   | `plugins`              | Dynamic `.so` plugin loading             |
| 10  | `app`                  | Application component management         |
| 11  | `scheduler`            | Cron-based task scheduling               |
| 12  | `network.ingress.http` | HTTP server (Gin), CORS, routing         |

Access a service from a controller:

```go
kernel := c.GetComponent().GetApp().GetKernel()
cfg := kernel.GetService("config")
```

## Configuration

The config service merges sources in order of precedence:

1. **Kernel defaults** — built-in `default.yaml`
2. **App config** — embedded YAML in app's `config/` directory
3. **Component configs** — embedded YAML per component
4. **Environment variables** — prefixed with `K2_` (e.g. `K2_HOST_INGRESS_PORT=8080` → `host.ingress.port`)

```go
cfg := kernel.GetService("config").(*config.Service)
cfg.Get("host.ingress.root")                    // any
cfg.GetAsString("host.ingress.root")             // (string, error)
cfg.GetAsInt("host.ingress.port")                // (int, error)
cfg.GetAsIntOrDefault("host.ingress.port", 8080) // (int, error)
cfg.Has("host.compute.enabled")                  // bool
cfg.SetValue("host.ingress.port", 9090)          // error
```

## Scheduling

Controllers can register cron schedules (6-field format with seconds):

```go
func (c *MyController) Start() error {
    c.AddSchedule("cleanup", "0 0 * * * *", func(fire string) string {
        // runs every hour
        return "done"
    })
    return nil
}
```

## Storage

### Object Stores

Two default stores are created: `root` (`{runDir}/`) and `local` (`{runDir}/home/`).

```go
stores := storeService.GetStores()         // []string
store, _ := storeService.GetStore("root")
exists, _ := store.Exists("path/to/object")
data, _ := store.ReadObject("path/to/object")
store.WriteObject("path/to/object", data)
store.DeleteObject("path/to/object")
```

Backends: **local** (filesystem) and **rclone** (cloud).

## Compute

The compute service manages infrastructure through platform providers:

| Provider     | Description                     |
| ------------ | ------------------------------- |
| `containers` | Docker/OCI container management |
| `kubernetes` | Kubernetes resource deployment  |
| `playbooks`  | Ansible playbook execution      |

Enabled via config key `host.compute.enabled`.

## Plugins

Dynamic `.so` plugins that export a `NewComponent` symbol:

```go
//go:build plugin

func NewComponent(app types.IApp) types.IAppComponent {
    return components.NewMyComponent(app)
}
```

Host app embeds plugins via `embed.FS` under `dist/`.

## Secrets & Profile

- **Profile** — persistent JSON at `{runDir}/etc/profile.json` with properties and secrets
- **Secrets** — base64-encoded values; auto-generates ECDSA P-256 SSH key pair on first init

## Project Structure

```
bases/              — base implementations (app, component, controller, service)
kernel/             — kernel runtime and all services
  app/              — application service
  compute/          — compute service + providers (containers, kubernetes, playbooks)
  config/           — configuration service
  monitoring/       — logging service
  network/ingress/  — HTTP ingress service + middlewares
  plugins/          — plugin loading service
  profile/          — profile service
  scheduler/        — scheduler service
  secrets/          — secrets service
  storage/          — paths, stores, volumes services
types/              — interfaces (IApp, IKernel, IKernelService, IAppComponent…)
system/             — utilities (exec, YAML/JSON marshaling, Go templates + Sprig)
specs/              — detailed specifications
samples/            — usage examples (simple, standalone, with_plugins, with-a-kube)
```

## Samples

| Sample         | Description                                            |
| -------------- | ------------------------------------------------------ |
| `simple`       | Minimal app with one component and one controller      |
| `standalone`   | App with compute (containers) and multiple controllers |
| `with_plugins` | Plugin host + dynamically loaded plugin                |
| `with-a-kube`  | App with Kubernetes integration                        |

## License

[GPL-3.0](LICENSE)
