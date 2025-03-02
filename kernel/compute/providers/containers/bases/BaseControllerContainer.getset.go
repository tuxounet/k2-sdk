package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
)

func (b *BaseControllerContainer) GetDefinition() *types.ContainerDefinition {
	return b.GetData("definition").(*types.ContainerDefinition)
}
