package scheduler

import (
	"github.com/robfig/cron/v3"
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = "scheduler"

type Service struct {
	bases.BaseKernelService
}

func NewService(k types.IKernel) types.IKernelService {

	base := bases.NewBaseKernelService(k, ServiceKey)

	service := &Service{base}

	cronInstance := cron.New(cron.WithSeconds())

	service.setCronRunner(cronInstance)
	return service

}
