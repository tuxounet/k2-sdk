package containers

import (
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
)

func (p *Provider) getPathsService() *paths.Service {
	return p.GetService().GetKernel().GetService(paths.ServiceKey).(*paths.Service)
}

func (s *Provider) getHostPortStart() (int, error) {

	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.containers.port.start")
	if err != nil {
		return -1, err
	}
	return port, nil

}
func (s *Provider) getHostPortEnd() (int, error) {
	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.containers.port.end")
	if err != nil {
		return -1, err
	}
	return port, nil

}
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
