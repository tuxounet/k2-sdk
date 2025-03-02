package ingress

import "github.com/gin-gonic/gin"

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
