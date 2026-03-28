# Spec 03 — Controllers

## Overview

Controllers handle HTTP routes and scheduled tasks within a component. They are the leaf nodes in the App → Component → Controller hierarchy.

## Controller Creation

```go
base := bases.NewBaseAppController(
    component,    // IAppComponent — parent component
    name,         // string — controller name (used in URL routing)
    order,        // int — execution priority
    config,       // *embed.FS — embedded config (optional)
    accessPolicy, // IAccessPolicy — "public" or "authenticated"
)
```

Returns a `BaseAppController` value (not pointer) — designed to be embedded in custom controller structs.

## Custom Controller Pattern

```go
type MyController struct {
    bases.BaseAppController
}

func NewMyController(component types.IAppComponent) types.IAppController {
    ctrl := &MyController{
        BaseAppController: bases.NewBaseAppController(component, "my-ctrl", 0, nil, types.AccessPolicyPublic),
    }
    return ctrl
}
```

## HTTP Route Registration

Override `Register(router *gin.RouterGroup)` to define routes:

```go
func (c *MyController) Register(router *gin.RouterGroup) error {
    router.GET("/hello", c.handleHello)
    router.POST("/data", c.handleData)
    return nil
}
```

Routes are mounted at: `/{baseUrl}/{component-name}/{controller-name}/`

## Lifecycle

| Method             | Description                                           |
| ------------------ | ----------------------------------------------------- |
| `Init()`           | No-op by default. Override for setup logic            |
| `Register(router)` | Receives a Gin router group scoped to this controller |
| `Start()`          | Called after all registration is complete             |
| `Stop()`           | Called during shutdown (reverse order)                |

## Data Store

Each controller has a key-value data store:

```go
controller.SetData("key", value)
val := controller.GetData("key") // returns any
```

- Returns `nil` for missing keys
- Scoped to the controller instance

## Schedule Attachment

Controllers can register cron schedules:

```go
func (c *MyController) Start() error {
    c.AddSchedule("my-task", "0 */5 * * * *", func(fire string) string {
        // runs every 5 minutes
        return "done"
    })
    return nil
}
```

See [05-scheduling.md](05-scheduling.md) for details.

## Accessing Kernel Services

```go
kernel := c.GetComponent().GetApp().GetKernel()
configService := kernel.GetService("config").(*config.Service)
value := configService.Get("my.setting")
```

## Sample Usage (from samples/simple)

```go
func NewSampleController(component types.IAppComponent) types.IAppController {
    return &SampleController{
        BaseAppController: bases.NewBaseAppController(component, "sample", 0, nil, types.AccessPolicyPublic),
    }
}

func (s *SampleController) Register(router *gin.RouterGroup) error {
    router.GET("/list", s.handleList)
    return nil
}
```
