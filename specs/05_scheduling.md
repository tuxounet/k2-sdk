# Spec 05 — Scheduling

## Overview

The scheduler service discovers cron schedules registered on controllers and executes them using `robfig/cron/v3`.

## Service Key

`"scheduler"` — accessible via `kernel.GetService("scheduler")`

## Defining a Schedule

Schedules are added to controllers via `AddSchedule`:

```go
func (c *MyController) Start() error {
    c.AddSchedule("cleanup", "0 0 * * * *", func(fire string) string {
        // runs every hour at minute 0
        return "cleaned up"
    })
    return nil
}
```

## Schedule Definition

```go
type IAppSchedule interface {
    GetName() string                    // schedule name
    GetCron() string                    // cron expression (6 fields, with seconds)
    GetTaskHandler() AppScheduleHandler // execution callback
}

type AppScheduleHandler = func(fire string) string
```

- `fire` parameter: timestamp string of when the job fired
- Return value: result string (logged by scheduler)

## Cron Expression Format

Uses **6-field cron expressions** (with seconds):

```
┌──────── second (0-59)
│ ┌────── minute (0-59)
│ │ ┌──── hour (0-23)
│ │ │ ┌── day of month (1-31)
│ │ │ │ ┌ month (1-12)
│ │ │ │ │ ┌ day of week (0-6)
│ │ │ │ │ │
* * * * * *
```

## Scheduler Lifecycle

1. **Start()** — Iterates App → Components → Controllers → Schedules
   - Registers each schedule with the cron runner via `AddFunc`
   - Starts the cron runner
2. **Stop()** — Stops the cron runner

## Execution Flow

For each schedule fire:

1. Log fire time
2. Call `schedule.GetTaskHandler()(fireTime)`
3. Log duration and result

## Sample Usage (from samples/simple)

```go
func (s *SampleController) Start() error {
    s.AddSchedule("hello", "0 */1 * * * *", func(fire string) string {
        return "world"
    })
    return nil
}
```
