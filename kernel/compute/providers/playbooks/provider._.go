package playbooks

import (
	computeBases "github.com/tuxounet/k2-sdk/kernel/compute/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"

	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Provider struct {
	computeBases.BasePlateformProvider[types.PlaybookDefinition]
}

const ProviderKey string = "playbooks"

func NewProvider(service runtimeTypes.IKernelService, ingressRegistar ingressTypes.IngressRegistarFunction) computeTypes.IPlateformProvider[types.PlaybookDefinition] {
	base := computeBases.NewBasePlateformProvider[types.PlaybookDefinition](service, ProviderKey, ingressRegistar)
	instance := &Provider{base}

	return instance
}
