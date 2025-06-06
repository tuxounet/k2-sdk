package stores

import (
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Service) GetBackends() ([]types.IStoreProvider, error) {
	data := s.GetData("backends")
	if data == nil {
		return nil, nil
	}
	return data.([]types.IStoreProvider), nil

}

func (s *Service) setBackends(backends []types.IStoreProvider) {
	s.SetData("backends", backends)
}

func (s *Service) setStores(stores []*types.Store) {
	s.SetData("stores", stores)
}

func (s *Service) getStores() []*types.Store {
	data := s.GetData("stores")
	if data == nil {
		return make([]*types.Store, 0)
	}
	return data.([]*types.Store)
}
