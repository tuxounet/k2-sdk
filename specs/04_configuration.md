# Spec 04 — Configuration

## Overview

The config service manages hierarchical YAML-based configuration with environment variable overrides, dot-notation access, and Go template rendering.

## Service Key

`"config"` — accessible via `kernel.GetService("config")`

## Configuration Sources (in order of precedence)

1. **Kernel defaults** — embedded in `kernel/config/defaults/default.yaml`
2. **App config** — embedded YAML files in the app's `config/` directory
3. **Component configs** — embedded YAML files per component
4. **Environment variables** — prefixed with `K2_`

Later sources override earlier ones (merge with override).

## Loading from Embedded FS

```go
configService.LoadFromEmbedFS("source-name", "folder", &myEmbedFS)
```

- Reads all `.yaml` files from the specified folder in the embed.FS
- Parses YAML and merges into the current config (with override)
- Silently skips if directory doesn't exist
- Skips non-YAML files and subdirectories

## Loading from Environment Variables

```go
configService.LoadFromEnvVars("source-name")
```

- Scans `os.Environ()` for variables prefixed with `K2_`
- Normalizes key: strips `K2_` prefix, lowercases, replaces `_` with `.`
- Example: `K2_HOST_INGRESS_PORT=8080` → `host.ingress.port = "8080"`
- Values are always strings from environment

## Dot-Notation Access

```go
configService.Get("host.ingress.root")        // any (nil if not found)
configService.Has("host.ingress.root")         // bool
configService.GetAsString("host.ingress.root") // (string, error)
configService.GetAsInt("host.ingress.port")    // (int, error) — returns -1 on error
configService.GetAsBool("host.compute.enabled") // (bool, error)
configService.GetAsStringOrDefault("key", "default") // (string, error)
configService.GetAsIntOrDefault("key", 8080)         // (int, error)
```

- Key traversal is **case-insensitive**
- Returns nil/error if path is not found or type doesn't match

## Setting Values

```go
configService.SetValue("host.ingress.port", 8080)
```

- Creates intermediate maps if they don't exist (auto-vivification)
- Returns error if key is empty or path is not traversable

## Template Rendering

```go
rendered, err := configService.Untemplate(templateBytes, "/api/v1")
```

Template data contains:

- `.Config` — full config map
- `.Name` — app name
- `.Version` — app version
- `.RootUrl` — value of `host.ingress.root`
- `.BasePath` — the base route (with trailing `/`)
- `.UIBasePath` — base route + `ui/`

Supports [Sprig](https://masterminds.github.io/sprig/) template functions.

## Sample Usage

```yaml
# config/app.yaml
app:
  greeting: "Hello World"
  max_items: 100
```

```go
kernel := c.GetComponent().GetApp().GetKernel()
cfg := kernel.GetService("config").(*config.Service)
greeting, _ := cfg.GetAsString("app.greeting")
```
