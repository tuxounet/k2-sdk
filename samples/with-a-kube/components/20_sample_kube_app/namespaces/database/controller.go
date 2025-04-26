package database

import (
	kubernetesBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/bases"
	kubernetesTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	kubernetesBases.BaseControllerKubernetes
}

func NewController(component types.IAppComponent) types.IAppController {
	base := kubernetesBases.NewBaseControllerKubernetes(component, &kubernetesTypes.NamespaceDefinition{
		Order:     10,
		Name:      "database",
		Templates: nil,
	}, nil)
	return &Controller{
		base,
	}
}
