package ingress

import (
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
)

func (s *Service) getConfigService() *config.Service {
	return s.GetKernel().GetService(config.ServiceKey).(*config.Service)
}

func (s *Service) getPathsService() *paths.Service {
	return s.GetKernel().GetService(paths.ServiceKey).(*paths.Service)
}
