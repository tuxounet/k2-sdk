package bases

import (
	"embed"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"

	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type BaseControllerContainer struct {
	runtimeBases.BaseAppController
}

func NewBaseControllerContainer(component runtimeTypes.IAppComponent, definition *types.ContainerDefinition, config *embed.FS) BaseControllerContainer {
	name := definition.Name
	base := runtimeBases.NewBaseAppController(component, name, definition.Order, config)
	instance := BaseControllerContainer{base}
	instance.SetData("definition", definition)

	instance.GetLogger().DebugF("BaseControllerContainer: [%d] %s", definition.Order, name)

	return instance
}
