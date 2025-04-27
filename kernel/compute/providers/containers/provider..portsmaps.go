package containers

// func (s *Provider) AllocateLocalPort(name string, order int, ingress *ingressTypes.IngressDefinition) (int, error) {

// 	portsmapStore := s.getPortsMapSore()
// 	records, err := portsmapStore.GetValue()
// 	if err != nil {
// 		return -1, err
// 	}
// 	allRecords := *records

// 	portStart, err := s.getHostPortStart()
// 	if err != nil {
// 		return -1, err
// 	}
// 	index := len(allRecords)

// 	localPort := portStart + index

// 	portEnd, err := s.getHostPortEnd()
// 	if err != nil {
// 		return -1, err
// 	}

// 	if localPort > portEnd {
// 		return -1, errors.New("no more ports available")
// 	}

// 	ingress.ServicePort = localPort

// 	record := containersTypes.PortsMapRecord{
// 		LocalPort:     localPort,
// 		ContainerName: name,
// 		Order:         order,
// 		Ingress:       ingress,
// 	}

// 	allRecords = append(allRecords, record)

// 	err = portsmapStore.SetValue(allRecords)
// 	if err != nil {
// 		s.GetLogger().ErrorF("Failed to write portmaps: %s", err)
// 		return -1, err
// 	}

// 	return record.LocalPort, nil
// }

// func (s *Service) LookupIngressPort(name string, order int, target string) (*types.PortsMapRecord, error) {
// 	portsmapStore := s.getPortsMapSore()
// 	records, err := portsmapStore.GetValue()
// 	if err != nil {
// 		return nil, err
// 	}
// 	allRecords := *records

// 	for _, port := range allRecords {
// 		if port.Order == order && port.ContainerName == name && port.Ingress != nil && port.Ingress.Path == target {
// 			return &port, nil
// 		}
// 	}

// 	s.GetLogger().Warn("LookupPort: no port found")

// 	return nil, nil

// }
