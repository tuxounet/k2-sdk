package engines

import (
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores"
	storeTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *PodmanEngine) getLocalHostAddress() (string, error) {

	kernel := s.service.GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	hostAddr, err := configService.GetAsString("host.address")
	if err != nil {
		return "", err
	}
	return hostAddr, nil

}

func (p *PodmanEngine) getPathsService() *paths.Service {
	return p.service.GetKernel().GetService(paths.ServiceKey).(*paths.Service)
}

func (p *PodmanEngine) getRootStore() (*storeTypes.Store, error) {
	stores := p.service.GetKernel().GetService(stores.ServiceKey).(*stores.Service)
	store, err := stores.GetStore("root")
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *PodmanEngine) getHostPortStart() (int, error) {

	kernel := s.service.GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.containers.port.start")
	if err != nil {
		return -1, err
	}
	return port, nil

}
func (s *PodmanEngine) getHostPortEnd() (int, error) {
	kernel := s.service.GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.containers.port.end")
	if err != nil {
		return -1, err
	}
	return port, nil
}
