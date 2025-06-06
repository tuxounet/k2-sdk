package containers

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	"github.com/tuxounet/k2-sdk/kernel/config"
	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Provider) setContainerEngine(engine types.IContainerEngine) {
	s.SetData("engine", engine)
}
func (s *Provider) getContainerEngine() types.IContainerEngine {
	engine, ok := s.GetData("engine").(types.IContainerEngine)
	if !ok {
		return nil
	}
	return engine
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
func (s *Provider) getPortsMapSore() storesTypes.IBaseObjectStore[[]types.PortsMapRecord] {
	return s.GetData("portsmap").(storesTypes.IBaseObjectStore[[]types.PortsMapRecord])

}
