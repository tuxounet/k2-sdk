package kubernetes

import (
	computeBases "github.com/tuxounet/k2-sdk/kernel/compute/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"

	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Provider struct {
	computeBases.BasePlateformProvider[types.NamespaceDefinition]
}

const ProviderKey string = "kubernetes"

func NewProvider(service runtimeTypes.IKernelService) computeTypes.IPlateformProvider[types.NamespaceDefinition] {

	base := computeBases.NewBasePlateformProvider[types.NamespaceDefinition](service, ProviderKey)
	instance := &Provider{base}

	return instance
}
