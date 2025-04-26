package admin

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
		Order:     5,
		Name:      "admin",
		Templates: &templates,
	}, nil)
	return &Controller{
		base,
	}
}
