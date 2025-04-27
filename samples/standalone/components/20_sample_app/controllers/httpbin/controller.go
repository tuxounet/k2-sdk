package httpbin

import (
	containersBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/bases"
	containersTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	containersBases.BaseControllerContainer
}

const ControllerKey = "httpbin"

func NewController(component types.IAppComponent) types.IAppController {

	base := containersBases.NewBaseControllerContainer(component, &containersTypes.ContainerDefinition{
		Order: 40,
		Name:  ControllerKey,
		Image: "docker.io/kennethreitz/httpbin:latest",
		Ingresses: []*containersTypes.ContainerDefinitionIngress{
			{
				AccessPolicy:  types.AccessPolicyAuthenticated,
				ContainerPort: 80,
				Path:          "/" + ControllerKey,
			},
		},
		Env: map[string]string{
			"TZ": "Europe/Paris",
		},
	}, nil)
	return &Controller{
		base,
	}
}
