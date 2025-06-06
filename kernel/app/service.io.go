package app

import (
	"github.com/tuxounet/k2-sdk/kernel/config"
)

func (s *Service) getConfigService() *config.Service {
	return s.GetKernel().GetService(config.ServiceKey).(*config.Service)
}
