package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes"
)

func (b *BaseControllerKubernetes) getComputeKubernetesProvider() *kubernetes.Provider {

	computeService := b.getComputeService()
	provider := computeService.GetProvider(kubernetes.ProviderKey)
	return provider.(*kubernetes.Provider)

}

func (b *BaseControllerKubernetes) getComputeService() *compute.Service {
	return b.GetComponent().GetApp().GetKernel().GetService(compute.ServiceKey).(*compute.Service)
}
