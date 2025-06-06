package ingress

import (
	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
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

func (s *Service) getIngressesRecords() []types.IngressDefinition {
	data := s.GetData("ingresses")
	if data == nil {
		return make([]types.IngressDefinition, 0)

	}
	return data.([]types.IngressDefinition)

}

func (s *Service) setIngressesRecords(ingresses []types.IngressDefinition) {
	s.SetData("ingresses", ingresses)
}
