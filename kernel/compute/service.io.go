package compute

import (
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores"
	storeTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Service) getPathsService() *paths.Service {
	return s.GetKernel().GetService(paths.ServiceKey).(*paths.Service)
}

func (s *Service) getConfigService() *config.Service {
	return s.GetKernel().GetService(config.ServiceKey).(*config.Service)
}

func (s *Service) getRootStore() (*storeTypes.Store, error) {
	stores := s.GetKernel().GetService(stores.ServiceKey).(*stores.Service)
	store, err := stores.GetStore("root")
	if err != nil {
		return nil, err
	}
	return store, nil
}
