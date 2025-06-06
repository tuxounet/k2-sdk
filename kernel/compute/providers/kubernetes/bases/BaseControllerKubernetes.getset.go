package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
)

func (b *BaseControllerKubernetes) GetDefinition() *types.NamespaceDefinition {
	return b.GetData("definition").(*types.NamespaceDefinition)
}
