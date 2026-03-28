# Spec 08 — Storage

## Overview

The storage layer provides three services: paths (utilities), stores (object persistence), and volumes (stub).

## Paths Service

**Key:** `"storage.paths"`

Utility methods for path manipulation:

| Method                 | Description                        |
| ---------------------- | ---------------------------------- |
| `CominePath(parts...)` | Joins path parts using `path.Join` |
| `SplitPath(parts...)`  | Joins then splits by `/`           |
| `GetBaseName(path)`    | Returns `path.Base`                |
| `GetDirName(path)`     | Returns `path.Dir`                 |
| `GetExtName(path)`     | Returns `path.Ext`                 |

## Stores Service

**Key:** `"storage.stores"`

Provides object storage with pluggable backends.

### Default Stores

| Name    | Backend | Path             |
| ------- | ------- | ---------------- |
| `root`  | local   | `{runDir}/`      |
| `local` | local   | `{runDir}/home/` |

### Store API

```go
stores := storeService.GetStores()        // []string
store, _ := storeService.GetStore("root")  // *Store
storeService.UpsertStore(store)            // create or update
```

### Store Operations

```go
exists, _ := store.Exists("path/to/object")
data, _ := store.ReadObject("path/to/object")
store.WriteObject("path/to/object", data)
store.DeleteObject("path/to/object")
```

### Backends

- **local** — filesystem-based storage
- **rclone** — cloud storage via rclone

## Volumes Service

**Key:** `"storage.volumes"`

Currently a stub — no methods implemented.

## Init Behavior

On initialization, the stores service:

1. Creates backend providers (local, rclone)
2. Sets up default stores (root, local)
3. Cleans `tmp/` directory and creates `.keep` marker
