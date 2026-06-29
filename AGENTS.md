# k2-sdk — Go SDK for kernel-driven applications

## Overview

**k2-sdk** is a Go SDK for building **kernel-driven applications** — apps with a structured lifecycle, modular components, and pluggable infrastructure services.

It provides a **host runtime** that manages the full lifecycle of a Go application through a layered architecture:

```
Application
├── Components (organizational units)
│   └── Controllers (HTTP routes & scheduled tasks)
├── Kernel
│   └── Services (infrastructure: database, messaging, secrets...)
└── Entrypoint
```

## Stack

- **Language:** Go 1.24+ (`go 1.24.0`)
- **Module:** `github.com/tuxounet/k2-sdk`
- **Entry point:** `entrypoint.go`
- **Build:** `make build` (check `Makefile`)

## Directory layout

| Path | Purpose |
|---|---|
| `entrypoint.go` | Application entry point / wiring |
| `bases/` | Base types and interfaces |
| `kernel/` | Core runtime — lifecycle, services, components |
| `system/` | System-level utilities |
| `types/` | Shared type definitions |
| `tools/` | Build/dev tools |
| `testutils/` | Test helpers and mocks |
| `specs/` | Specifications |
| `samples/` | Example applications |
| `Makefile` | Build / test targets |
| `version.txt` | Current version |

## Key concepts

- **Components** — organizational units that group controllers and sub-components.
- **Controllers** — handle HTTP routes, scheduled tasks, and lifecycle events.
- **Kernel services** — interchangeable backends (database, messaging, secrets, etc.) wired at startup.
- **Lifecycle** — the host runtime manages init, start, ready, stop phases.

## Build & test

```sh
make build     # build SDK library
make test      # run tests
```

## Dependencies

Check `go.mod` for full dependency list. Requires Go 1.24+.

## Related projects

- **k2** (`../k2/`) — The YAML template engine; some scaffolding uses k2
- **k-assets** (`../k-assets/`) — Built on k2-sdk
- **k-office** (`../k-office/`) — Built on k2-sdk
- **k-kube** (`../k-kube/`) — Built on k2-sdk