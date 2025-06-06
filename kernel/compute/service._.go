package compute

import (
	runtimeBases "github.com/tuxounet/k2-sdk/bases"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = "compute"

type Service struct {
	runtimeBases.BaseKernelService
}

func NewService(k runtimeTypes.IKernel) runtimeTypes.IKernelService {

	base := runtimeBases.NewBaseKernelService(k, ServiceKey)
	instance := &Service{base}

	return instance
}
