package types

type IAppSchedule interface {
	GetName() string
	GetCron() string
	GetTaskHandler() AppScheduleHandler
}

type AppScheduleHandler = func(fire string) string
