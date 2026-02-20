package bases

import (
	"embed"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"

	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type BaseControllerKubernetes struct {
	runtimeBases.BaseAppController
}

func NewBaseControllerKubernetes(component runtimeTypes.IAppComponent, definition *types.NamespaceDefinition, config *embed.FS) BaseControllerKubernetes {
	name := definition.Name
	base := runtimeBases.NewBaseAppController(component, name, definition.Order, config, runtimeTypes.AccessPolicyPublic)
	instance := BaseControllerKubernetes{base}
	instance.SetData("definition", definition)

	instance.GetLogger().DebugF("NewBaseControllerKubernetes: [%d] %s", definition.Order, name)

	return instance
}
