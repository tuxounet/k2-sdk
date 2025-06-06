package containers

import (
	"path/filepath"

	"github.com/tuxounet/k2-sdk/kernel/compute/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"

	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	storesBases "github.com/tuxounet/k2-sdk/kernel/storage/stores/bases"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Provider struct {
	bases.BasePlateformProvider[types.ContainerDefinition]
}

const ProviderKey string = "containers"

func NewProvider(service runtimeTypes.IKernelService, ingressRegistar ingressTypes.IngressRegistarFunction) computeTypes.IPlateformProvider[types.ContainerDefinition] {
	base := bases.NewBasePlateformProvider[types.ContainerDefinition](service, ProviderKey, ingressRegistar)
	instance := &Provider{base}

	portsMapStore := storesBases.NewObjectStore[[]types.PortsMapRecord](
		service.GetKernel(),
		instance, "root",
		filepath.Join("etc", "compute", "containers", "portsmap.json"),
		"[]",
	)

	instance.SetData("portsmap", portsMapStore)
	return instance
}
