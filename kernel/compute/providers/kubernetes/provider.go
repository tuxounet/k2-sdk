package kubernetes

import (
	computeBases "github.com/tuxounet/k2-sdk/kernel/compute/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"

	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Provider struct {
	computeBases.BasePlateformProvider[types.NamespaceDefinition]
}

const ProviderKey string = "kubernetes"

func NewProvider(service runtimeTypes.IKernelService, ingressRegistar ingressTypes.IngressRegistarFunction) computeTypes.IPlateformProvider[types.NamespaceDefinition] {
	base := computeBases.NewBasePlateformProvider[types.NamespaceDefinition](service, ProviderKey, ingressRegistar)
	instance := &Provider{base}

	return instance
}
