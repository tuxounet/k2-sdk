package loadbalancer

import (
	_ "embed"

	playbooksBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/bases"
	playbooksTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	playbooksBases.BaseControllerPlaybook
}

//go:embed tasks/provision.yaml
var provisionTasks string

//go:embed tasks/teardown.yaml
var teardownTasks string

const ControllerKey = "loadbalancer"

func NewController(component types.IAppComponent) types.IAppController {

	base := playbooksBases.NewBaseControllerPlaybook(component, &playbooksTypes.PlaybookDefinition{
		Order:     50,
		Name:      ControllerKey,
		Provision: provisionTasks,
		Teardown:  teardownTasks,
	}, nil)
	return &Controller{
		base,
	}
}
