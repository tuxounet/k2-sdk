package browser

import (
	"fmt"

	containersBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/bases"
	containersTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	containersBases.BaseControllerContainer
}

const ControllerKey = "browser"

func NewController(component types.IAppComponent) types.IAppController {

	base := containersBases.NewBaseControllerContainer(component, &containersTypes.ContainerDefinition{
		Order: 40,
		Name:  ControllerKey,
		Image: "docker.io/filebrowser/filebrowser:v2.31.2",
		Ingresses: []*containersTypes.ContainerDefinitionIngress{
			{
				AccessPolicy:  types.AccessPolicyAdmin,
				ContainerPort: 8080,
				Path:          fmt.Sprintf("/%s/%s/", component.GetName(), ControllerKey),
			},
		},
		Env: map[string]string{
			"TZ":          "Europe/Paris",
			"FB_NOAUTH":   "true",
			"LC_ALL":      "fr_FR.UTF-8",
			"FB_PORT":     "8080",
			"FB_ROOT":     "/srv",
			"FB_LOG":      "stdout",
			"FB_DATABASE": "/config/filebrowser.db",
			"FB_BASEURL":  fmt.Sprintf("/%s/%s/", component.GetName(), ControllerKey),
		},
		Volumes: []*containersTypes.ContainerDefinitionVolume{
			{
				Name:          "srv",
				ContainerPath: "/srv",
				Binding: containersTypes.ContainerDefinitionVolumeBinding{
					Type:     containersTypes.ContainerDefinitionVolumeBindingTypeMount,
					HostPath: component.GetApp().GetKernel().GetRunDirectory(),
				},
			},
			{
				Name:          "config",
				ContainerPath: "/config",
				Binding: containersTypes.ContainerDefinitionVolumeBinding{
					Type: containersTypes.ContainerDefinitionVolumeBindingTypeMount,
				},
			},
		},
	}, nil)
	return &Controller{
		base,
	}
}
