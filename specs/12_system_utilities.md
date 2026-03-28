# Spec 12 — System Utilities

## Overview

The `system` package provides low-level utilities for command execution, data marshaling, and Go template rendering.

## Command Execution

### CmdCall

```go
cmd := system.NewCmdCall(logger, "ansible-playbook", "-i", "hosts.yaml", "playbook.yaml")
cmd.Cwd = &workDir  // optional working directory
fmt.Println(cmd.String()) // "ansible-playbook -i hosts.yaml playbook.yaml"
```

### Execution Functions

| Function                  | Description                            |
| ------------------------- | -------------------------------------- |
| `OsExec(cmd)`             | Execute and fail on non-zero exit code |
| `OsExecWithExitCode(cmd)` | Execute and return exit code           |
| `OsExecAndTailToLog(cmd)` | Execute with output streamed to logger |

- All functions log the command at trace level
- `OsExec` returns error if exit code ≠ 0
- `OsExecAndTailToLog` captures stdout/stderr in separate log streams

## Data Marshaling

### YAML

```go
result, err := system.LoadYamlFromString[MyStruct](yamlString)
yamlStr, err := system.DumpToYamlString(data)
```

### JSON

```go
result, err := system.LoadJSONFromString[MyStruct](jsonString)
jsonStr, err := system.DumpToJsonString(data)
```

- Generic functions using Go 1.18+ type parameters
- Returns zero value + error on parse failure

## Template Rendering

```go
result, err := system.UnTemplateWithGoTemplate(templateStr, data)
```

- Uses Go's `text/template` package
- Includes [Sprig](https://masterminds.github.io/sprig/) template functions
- Template name is auto-generated from MD5 hash of template content
- Returns error on parse or execution failure

### Template Example

```go
tpl := "Hello {{ .Name }}, you have {{ .Count | add 1 }} items"
data := map[string]any{"Name": "World", "Count": 5}
result, _ := system.UnTemplateWithGoTemplate(tpl, data)
// result: "Hello World, you have 6 items"
```
