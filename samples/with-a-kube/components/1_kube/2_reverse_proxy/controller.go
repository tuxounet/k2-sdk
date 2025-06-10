package reverse_proxy

import (
	_ "embed"

	containersBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/bases"
	containersTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	containersBases.BaseControllerContainer
}

const ControllerKey = "reverse_proxy"

//go:embed files/haproxy.cfg
var haproxyConfig string

func NewController(component types.IAppComponent) types.IAppController {

	base := containersBases.NewBaseControllerContainer(component, &containersTypes.ContainerDefinition{
		Order: 51,
		Name:  ControllerKey,
		Image: "docker.io/haproxy:lts",

		Env: map[string]string{
			"TZ":     "Europe/Paris",
			"LC_ALL": "fr_FR.UTF-8",
		},
		Volumes: []*containersTypes.ContainerDefinitionVolume{
			{
				Name:          "config",
				ContainerPath: "/usr/local/etc/haproxy/haproxy.cfg",
				Binding: containersTypes.ContainerDefinitionVolumeBinding{
					Type:    containersTypes.ContainerDefinitionVolumeBindingTypeContent,
					Content: haproxyConfig,
				},
			},
		},
		Ports: []*containersTypes.ContainerDefinitionPort{
			{
				ContainerPort: "80",
				HostPort:      "61080",
				Protocol:      "tcp",
			},
		},
	}, nil)
	return &Controller{
		base,
	}
}
