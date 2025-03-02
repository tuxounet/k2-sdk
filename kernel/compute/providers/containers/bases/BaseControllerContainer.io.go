package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers"
	"github.com/tuxounet/k2-sdk/kernel/config"
)

func (b *BaseControllerContainer) getComputeContainersProviders() *containers.Provider {

	computeService := b.getComputeService()
	provider := computeService.GetProvider(containers.ProviderKey)
	return provider.(*containers.Provider)

}

func (b *BaseControllerContainer) getComputeService() *compute.Service {
	return b.GetComponent().GetApp().GetKernel().GetService(compute.ServiceKey).(*compute.Service)
}

func (b *BaseControllerContainer) getConfigService() *config.Service {
	return b.GetComponent().GetApp().GetKernel().GetService(config.ServiceKey).(*config.Service)
}
