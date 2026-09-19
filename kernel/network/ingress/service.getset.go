package ingress

import (
	"net/http"

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

// getHTTPServer returns the stored *http.Server for the main HTTP/TLS server.
// This is set during Listen() and used by Stop() to perform graceful shutdown.
func (s *Service) getHTTPServer() *http.Server {
	data := s.GetData("httpServer")
	if data == nil {
		return nil
	}
	return data.(*http.Server)
}

// setHTTPServer stores the *http.Server for later shutdown.
func (s *Service) setHTTPServer(server *http.Server) {
	s.SetData("httpServer", server)
}

// getRedirectServer returns the stored *http.Server for the HTTP→HTTPS redirect server.
func (s *Service) getRedirectServer() *http.Server {
	data := s.GetData("redirectServer")
	if data == nil {
		return nil
	}
	return data.(*http.Server)
}

// setRedirectServer stores the redirect *http.Server for later shutdown.
func (s *Service) setRedirectServer(server *http.Server) {
	s.SetData("redirectServer", server)
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
