package ingress

import (
	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Service) GetServer() *gin.Engine {
	data := s.GetData("server")
	if data == nil {
		s.GetLogger().Warn("Server not found")
		return nil
	}

	return data.(*gin.Engine)
}

func (s *Service) setServer(server *gin.Engine) {
	s.SetData("server", server)
}

func (s *Service) GetRouter() *gin.RouterGroup {
	data := s.GetData("router")
	if data == nil {
		s.GetLogger().Warn("Router not found")
		return nil
	}

	return data.(*gin.RouterGroup)
}

func (s *Service) getIngressesStore() storesTypes.IBaseObjectStore[[]types.IngressDefinition] {
	return s.GetData("ingresses").(storesTypes.IBaseObjectStore[[]types.IngressDefinition])

}
