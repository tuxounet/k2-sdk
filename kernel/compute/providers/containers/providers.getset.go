package containers

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	"github.com/tuxounet/k2-sdk/kernel/config"
	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Provider) getLocalHostAddress() (string, error) {

	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	hostAddr, err := configService.GetAsString("host.address")
	if err != nil {
		return "", err
	}
	return hostAddr, nil

}

func (s *Provider) getPortsMapSore() storesTypes.IBaseObjectStore[[]types.PortsMapRecord] {
	return s.GetData("portsmap").(storesTypes.IBaseObjectStore[[]types.PortsMapRecord])

}
