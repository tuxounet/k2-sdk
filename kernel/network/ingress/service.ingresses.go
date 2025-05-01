package ingress

import (
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
)

func (s *Service) RegisterIngress(ingress *types.IngressDefinition) error {
	records := s.getIngressesRecords()
	records = append(records, *ingress)
	s.setIngressesRecords(records)
	s.GetLogger().InfoF("Ingress %s registered", ingress.IngressPath)
	return nil

}
