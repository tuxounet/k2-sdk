package kubernetes

import (
	"path/filepath"

	computeBases "github.com/tuxounet/k2-sdk/kernel/compute/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	storesBases "github.com/tuxounet/k2-sdk/kernel/storage/stores/bases"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Provider struct {
	computeBases.BasePlateformProvider[types.NamespaceDefinition]
}

const ProviderKey string = "kubernetes"

func NewProvider(service runtimeTypes.IKernelService, ingressRegistar ingressTypes.IngressRegistarFunction) computeTypes.IPlateformProvider[types.NamespaceDefinition] {
	base := computeBases.NewBasePlateformProvider[types.NamespaceDefinition](service, ProviderKey, ingressRegistar)
	instance := &Provider{base}
	portsForwardStore := storesBases.NewObjectStore[[]types.PortsForwardRecord](
		service.GetKernel(),
		instance, "root",
		filepath.Join("etc", "compute", "kubernetes", "forwards.json"),
		"[]",
	)
	instance.SetData("forwards", portsForwardStore)
	instance.setForwarders(make([]*types.PortForwarder, 0))
	return instance
}
