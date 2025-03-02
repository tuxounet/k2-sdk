package containers

import (
	"errors"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	"github.com/tuxounet/k2-sdk/kernel/config"
)

func (s *Provider) RegisterIngress(name string, order int, ingress *types.ContainerDefinitionIngress) (int, error) {

	portsmapStore := s.getPortsMapSore()
	records, err := portsmapStore.GetValue()
	if err != nil {
		return -1, err
	}
	allRecords := *records

	portStart, err := s.getHostPortStart()
	if err != nil {
		return -1, err
	}
	index := len(allRecords)

	localPort := portStart + index

	portEnd, err := s.getHostPortEnd()
	if err != nil {
		return -1, err
	}

	if localPort > portEnd {
		return -1, errors.New("no more ports available")
	}

	record := types.PortsMapRecord{
		LocalPort:     localPort,
		ContainerName: name,
		Order:         order,
		Ingress:       ingress,
	}

	allRecords = append(allRecords, record)

	err = portsmapStore.SetValue(allRecords)
	if err != nil {
		s.GetLogger().ErrorF("Failed to write portmaps: %s", err)
		return -1, err
	}

	return record.LocalPort, nil
}

func (s *Provider) LookupIngressPort(name string, order int, target string) (*types.PortsMapRecord, error) {
	portsmapStore := s.getPortsMapSore()
	records, err := portsmapStore.GetValue()
	if err != nil {
		return nil, err
	}
	allRecords := *records

	for _, port := range allRecords {
		if port.Order == order && port.ContainerName == name && port.Ingress != nil && port.Ingress.Path == target {
			return &port, nil
		}
	}

	s.GetLogger().Warn("LookupPort: no port found")

	return nil, nil

}

func (s *Provider) getHostPortStart() (int, error) {

	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.port.start")
	if err != nil {
		return -1, err
	}
	return port, nil

}
func (s *Provider) getHostPortEnd() (int, error) {
	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.port.end")
	if err != nil {
		return -1, err
	}
	return port, nil

}
