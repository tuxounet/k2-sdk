package tls_ca_issuer

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

const ControllerKey = "tls_ca_issuer"

func NewController(component types.IAppComponent) types.IAppController {

	base := playbooksBases.NewBaseControllerPlaybook(component, &playbooksTypes.PlaybookDefinition{
		Order:     60,
		Name:      ControllerKey,
		Provision: provisionTasks,
	}, &controllerConfig)
	return &Controller{
		base,
	}
}
