package bases

import "github.com/gin-gonic/gin"

func (b *BaseControllerPlaybook) Register(r *gin.RouterGroup) error {

	provider := b.getComputePlaybooksProvider()

	definition := b.GetDefinition()
	err := provider.RegisterDefinition(*definition)

	if err != nil {
		b.GetLogger().ErrorF("Failed to register container defintion inside provider: %s", err)
		return err
	}

	return nil
}
