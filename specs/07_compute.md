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
4. **Stop()** — Call `provider.Stop()` → execute stop playbook → teardown playbook

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
