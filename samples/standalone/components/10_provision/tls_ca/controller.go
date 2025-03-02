package tls_ca

import (
	"embed"

	playbooksBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/bases"
	playbooksTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	playbooksBases.BaseControllerPlaybook
}

//go:embed config/*.yaml
var controllerConfig embed.FS

//go:embed playbooks/provision.yaml
var provisionPlaybook string

const ControllerKey = "tls_ca"

func NewController(component types.IAppComponent) types.IAppController {

	base := playbooksBases.NewBaseControllerPlaybook(component, &playbooksTypes.PlaybookDefinition{
		Order:     3,
		Name:      ControllerKey,
		Provision: provisionPlaybook,
	}, &controllerConfig)
	return &Controller{
		base,
	}
}
