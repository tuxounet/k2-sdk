package ingress

import (
	"embed"

	playbooksBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/bases"
	playbooksTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	playbooksBases.BaseControllerPlaybook
}

//go:embed tasks/provision.yaml
var provisionTasks string

//go:embed config/*.yaml
var controllerConfig embed.FS

const ControllerKey = "ingress"

func NewController(component types.IAppComponent) types.IAppController {

	base := playbooksBases.NewBaseControllerPlaybook(component, &playbooksTypes.PlaybookDefinition{
		Order:     51,
		Name:      ControllerKey,
		Provision: provisionTasks,
	}, &controllerConfig)
	return &Controller{
		base,
	}
}
