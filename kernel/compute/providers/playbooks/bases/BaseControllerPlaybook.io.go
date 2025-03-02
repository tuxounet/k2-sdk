package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks"
)

func (b *BaseControllerPlaybook) getComputePlaybooksProvider() *playbooks.Provider {

	computeService := b.getComputeService()
	provider := computeService.GetProvider(playbooks.ProviderKey)
	return provider.(*playbooks.Provider)

}

func (b *BaseControllerPlaybook) getComputeService() *compute.Service {
	return b.GetComponent().GetApp().GetKernel().GetService(compute.ServiceKey).(*compute.Service)
}
