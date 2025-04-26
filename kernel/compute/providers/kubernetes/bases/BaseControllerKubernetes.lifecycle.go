package bases

import (
	"github.com/gin-gonic/gin"
)

func (b *BaseControllerKubernetes) Init() error {
	provider := b.getComputeKubernetesProvider()

	definition := b.GetDefinition()
	err := provider.RegisterDefinition(*definition)

	if err != nil {
		b.GetLogger().ErrorF("Failed to register definition: %s", err.Error())
		return err
	}

	return nil
}

func (b *BaseControllerKubernetes) Register(r *gin.RouterGroup) error {

	return nil
}
