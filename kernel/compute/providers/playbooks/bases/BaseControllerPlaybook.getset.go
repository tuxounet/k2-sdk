package bases

import "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"

func (b *BaseControllerPlaybook) GetDefinition() *types.PlaybookDefinition {
	return b.GetData("definition").(*types.PlaybookDefinition)
}
