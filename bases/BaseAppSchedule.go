package bases

import "github.com/tuxounet/k2-sdk/types"

type BaseAppSchedule struct {
	name        string
	cron        string
	taskHandler types.AppScheduleHandler
}

func NewBaseAppSchedule(name string, cron string, handler types.AppScheduleHandler) BaseAppSchedule {
	return BaseAppSchedule{
		name:        name,
		cron:        cron,
		taskHandler: handler,
	}
}

func (b *BaseAppSchedule) GetName() string {
	return b.name
}

func (b *BaseAppSchedule) GetCron() string {
	return b.cron
}

func (b *BaseAppSchedule) GetTaskHandler() types.AppScheduleHandler {
	return b.taskHandler
}
