package bases

import (
	"embed"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"

	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type BaseControllerPlaybook struct {
	runtimeBases.BaseAppController
}

func NewBaseControllerPlaybook(component runtimeTypes.IAppComponent, definition *types.PlaybookDefinition, config *embed.FS) BaseControllerPlaybook {
	name := definition.Name
	base := runtimeBases.NewBaseAppController(component, name, definition.Order, config, runtimeTypes.AccessPolicyPublic)
	instance := BaseControllerPlaybook{base}
	instance.SetData("definition", definition)

	return instance
}
