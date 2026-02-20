package plugins

import (
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = "plugins"

type PluginsService struct {
	bases.BaseKernelService
}

func NewService(k types.IKernel) types.IKernelService {

	base := bases.NewBaseKernelService(k, ServiceKey)
	return &PluginsService{base}

}
