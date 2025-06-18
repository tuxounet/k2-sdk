package tty

import (
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
		Order: 31,
		Name:  ControllerKey,
		Image: "docker.io/wettyoss/wetty:latest",
		Ingresses: []*containersTypes.ContainerDefinitionIngress{
			{
				AccessPolicy:  types.AccessPolicyPublic,
				ContainerPort: 3000,
				Path:          "/tty/",
			},
		},
		Env: map[string]string{
			"TZ": "Europe/Paris",
		},
		Command: &[]string{
			"pnpm",
			"start",
			"--host",
			"0.0.0.0",
			"--port",
			"3000",
			"--base",
			"/tty/shell",
			"--command",
			"/bin/sh",
		},
	}, nil)
	return &Controller{
		base,
	}
}
