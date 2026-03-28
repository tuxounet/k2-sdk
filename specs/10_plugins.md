# Spec 10 — Plugins

## Overview

The plugin system allows dynamic loading of Go shared object (`.so`) files that provide additional application components.

## Service Key

`"plugins"` — accessible via `kernel.GetService("plugins")`

## Plugin Contract

A plugin must:

1. Be compiled as a Go plugin (`-buildmode=plugin`)
2. Export a symbol named `NewComponent` with the signature:
   ```go
   type PluginComponentCtor = func(app IApp) IAppComponent
   ```
3. Use the build tag `//go:build plugin`

## Plugin Embedding

Plugins are embedded in the host application's `embed.FS` under a `dist/` directory:

```go
//go:embed dist/*.so
var pluginRegistry embed.FS
```

The embed.FS is passed as the `externalComponents` parameter to `NewBaseApp`.

## Loading Flow

1. Get app's externals (`GetExternals()`)
2. Scan `dist/` directory for `.so` files
3. Extract each `.so` to `{runDir}/var/lib/`
4. Load using Go's `plugin.Open()`
5. Lookup `NewComponent` symbol
6. Store as `PluginInfo{Name, Path}`

## PluginInfo

```go
type PluginInfo struct {
    Name string // plugin filename (without extension)
    Path string // extracted filesystem path
}
```

## Constants

- Plugin symbol name: `"NewComponent"`
- Plugin extension: `".so"`
- Plugin embed directory: `"dist"`
- Plugin extraction directory: `{runDir}/var/lib/`

## Sample Usage (from samples/with_plugins)

### Plugin (plug_hello):

```go
//go:build plugin

func NewComponent(app types.IApp) types.IAppComponent {
    return components.NewHelloComponent(app)
}
```

### Host (plug_host):

```go
//go:embed dist/*.so
var pluginRegistry embed.FS

func NewApp() types.IApp {
    return bases.NewBaseApp("host", "0.0.1", nil, nil, nil,
        []types.AppComponentCtor{}, &pluginRegistry)
}
```
