package ingress

import (
	"path/filepath"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	storesBases "github.com/tuxounet/k2-sdk/kernel/storage/stores/bases"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = "network.ingress.http"

type Service struct {
	runtimeBases.BaseKernelService
}

func NewService(k runtimeTypes.IKernel) runtimeTypes.IKernelService {

	base := runtimeBases.NewBaseKernelService(k, ServiceKey)
	instance := &Service{base}
	ingressesStore := storesBases.NewObjectStore[[]types.IngressDefinition](
		k,
		instance, "root",
		filepath.Join("etc", "network", "ingresses.json"),
		"[]",
	)

	instance.SetData("ingresses", ingressesStore)
	return instance
}
