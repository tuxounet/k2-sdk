package containers

import (
	"github.com/tuxounet/k2-sdk/kernel/config"
)

func (s *Provider) getEngine() string {
	defaultEngine := "podman"
	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	engine, err := configService.GetAsString("host.compute.containers.engine")
	if err != nil {
		s.GetLogger().ErrorF("Failed to get engine: %s", err)
		return defaultEngine
	}
	return engine

}
