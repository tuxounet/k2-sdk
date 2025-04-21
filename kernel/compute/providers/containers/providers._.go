package containers

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/engines"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"

	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	storesBases "github.com/tuxounet/k2-sdk/kernel/storage/stores/bases"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Provider struct {
	bases.BasePlateformProvider[types.ContainerDefinition]
}

const ProviderKey string = "containers"

func NewProvider(service runtimeTypes.IKernelService) computeTypes.IPlateformProvider[types.ContainerDefinition] {
	base := bases.NewBasePlateformProvider[types.ContainerDefinition](service, ProviderKey)
	instance := &Provider{base}

	paths := instance.getPathsService()

	portsMapStore := storesBases.NewObjectStore[[]types.PortsMapRecord](
		service.GetKernel(),
		instance, "root",
		paths.CominePath("etc", "compute", ProviderKey, "portsmap.json"),
		"[]",
	)

	instance.SetData("portsmap", portsMapStore)

	engine := instance.getEngine()
	switch engine {
	case "podman":
		instance.setContainerEngine(engines.NewPodmanEngine(service, instance.RegisterIngress))
	default:
		panic("engine not supported")
	}

	return instance
}
