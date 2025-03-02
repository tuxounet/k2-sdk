package compute

import (
	"github.com/tuxounet/k2-sdk/bases"

	"github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = "compute"

type Service struct {
	bases.BaseKernelService
}

func NewService(k types.IKernel) types.IKernelService {

	base := bases.NewBaseKernelService(k, ServiceKey)
	instance := &Service{base}

	return instance
}
