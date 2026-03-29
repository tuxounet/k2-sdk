# Spec 07 — Compute

## Overview

The compute service manages infrastructure provisioning through platform providers (containers, Kubernetes, playbooks). It uses an Ansible-based runner system to orchestrate lifecycle operations.

## Service Key

`"compute"` — accessible via `kernel.GetService("compute")`

## Platform Providers

Three provider types are available:

| Provider     | Description                     |
| ------------ | ------------------------------- |
| `containers` | Docker/OCI container management |
| `kubernetes` | Kubernetes resource deployment  |
| `playbooks`  | Ansible playbook execution      |

## Runner Definitions

Runners define infrastructure operations:

```go
type RunnerDefinition struct {
    Order     int    // execution priority
    Provider  string // "container", "kubernetes", "playbook"
    Name      string // unique runner name
    Provision string // provision task/playbook
    Teardown  string // teardown task/playbook
    Start     string // start task/playbook
    Stop      string // stop task/playbook
}
```

## Runner Verbs

```go
const (
    RunnerVerbProvision RunnerVerb = "provision"
    RunnerVerbStart     RunnerVerb = "start"
    RunnerVerbStop      RunnerVerb = "stop"
    RunnerVerbTeardown  RunnerVerb = "teardown"
)
```

## Lifecycle

1. **Init()** — Create platform providers, call `provider.Init()`
2. **Register()** — Render Ansible inventory, collect runners from providers, generate 4 playbooks (provision/start/stop/teardown)
3. **Start()** — Execute provision playbook → start playbook → call `provider.Start()`
4. **Stop()** — Call `provider.Stop()` → execute stop playbook
5. **ExecVerb(verb)** — Execute a specific verb playbook (used by `HostProvisionOnly`/`HostTeardownOnly`)

> **Note:** `teardown` is **never** called during normal `Stop()`. It is only executed explicitly via `HostTeardownOnly()` or `ExecVerb(teardown)`. This preserves provisioned resources across restarts.

## Checksum Cache

To optimize execution times, the compute service implements a checksum-based cache for playbook execution:

- **Scope**: Per-verb (provision, start, stop, teardown) — each verb's rendered playbook is checksummed independently
- **Input**: SHA-256 of the rendered playbook YAML content + serialized configuration map. This ensures both playbook changes and config value changes trigger re-execution
- **TTL**: Configurable via `host.compute.cache.ttl` (in minutes, default: 1440 = 24 hours) — playbooks are re-executed after TTL even without changes, to guard against infrastructure drift
- **Cache location**: `{runDir}/etc/compute/checksums.json`
- **Force flag**: `kernel.IsForceCompute()` bypasses the cache entirely (set via `--force-compute` CLI flag)

### Skip conditions

A playbook verb is **skipped** when all of these are true:

1. A cached checksum entry exists for this verb
2. The current checksum matches the cached checksum (no playbook or config changes)
3. Less than 24 hours have elapsed since the last execution

### Cache invalidation on lifecycle verbs

When a "down" verb executes, it invalidates the cache of its counterpart "up" verb so the next cycle re-executes from scratch:

| Verb executed | Cache invalidated                       | Rationale                                                            |
| ------------- | --------------------------------------- | -------------------------------------------------------------------- |
| `stop`        | `start`                                 | Infrastructure was stopped; `start` must re-run on next cycle        |
| `teardown`    | `provision` + **all** (full cache nuke) | Resources were deallocated; `provision` must re-run to recreate them |

This ensures that after a `stop`/`teardown`, the next `start`/`provision` always runs regardless of checksum state.

### Cache data structure

```go
type VerbChecksum struct {
    Verb       RunnerVerb `json:"verb"`
    Checksum   string     `json:"checksum"`
    ExecutedAt time.Time  `json:"executed_at"`
}
type ChecksumCache map[string]VerbChecksum
```

## Configuration

Enabled via config key: `host.compute.enabled` (bool)

## Inventory

Generates Ansible `hosts.yaml` and `group_vars/all.yaml` from config data.

## Sample Usage (from samples/standalone)

```go
type ProvisionController struct {
    compute.BaseControllerContainer
}

func (c *ProvisionController) Init() error {
    c.RegisterDefinition(containers.ContainerDefinition{
        Name:  "browser",
        Image: "jlesage/firefox:latest",
        Ports: []string{"5800:5800"},
    })
    return nil
}
```
