package database

import (
	"embed"

	kubernetesBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/bases"
	kubernetesTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	kubernetesBases.BaseControllerKubernetes
}

//go:embed **/*.yaml
var templates embed.FS

func NewController(component types.IAppComponent) types.IAppController {
	base := kubernetesBases.NewBaseControllerKubernetes(component, &kubernetesTypes.NamespaceDefinition{
		Order:     10,
		Name:      "database",
		Templates: &templates,
		// Ports: []kubernetesTypes.NamespacePortForwards{
		// 	{
		// 		LocalPort:        38000,
		// 		ServiceNamespace: "database",
		// 		ServiceName:      "db",
		// 		ServicePort:      5432,
		// 	},
		// },
	}, nil)
	return &Controller{
		base,
	}
}
