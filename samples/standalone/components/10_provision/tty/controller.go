package tty

import (
	"fmt"

	containersBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/bases"
	containersTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	containersBases.BaseControllerContainer
}

const ControllerKey = "tty"

func NewController(component types.IAppComponent) types.IAppController {

	base := containersBases.NewBaseControllerContainer(component, &containersTypes.ContainerDefinition{
		Order: 30,
		Name:  ControllerKey,
		Image: "docker.io/wettyoss/wetty:latest",
		Ingresses: []*containersTypes.ContainerDefinitionIngress{
			{
				AccessPolicy:  types.AccessPolicyPublic,
				ContainerPort: 3000,
				Path:          fmt.Sprintf("/%s/%s/", component.GetName(), ControllerKey),
			},
		},
		Env: map[string]string{
			"TZ": "Europe/Paris",
		},
		Command: &[]string{
			"pnpm",
			"start",
			"--port",
			"3000",
			"--base",
			"/provision/tty/",
			"--command",
			"bash",
		},
	}, nil)
	return &Controller{
		base,
	}
}
