package ingress

import (
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
)

func (s *Service) RegisterIngress(ingress *types.IngressDefinition) error {
	store := s.getIngressesStore()

	records, err := store.GetValue()
	if err != nil {
		return err
	}
	allRecords := *records
	allRecords = append(allRecords, *ingress)
	err = store.SetValue(allRecords)
	if err != nil {
		return err
	}
	s.GetLogger().InfoF("Ingress %s registered", ingress.IngressHost)
	return nil
}
