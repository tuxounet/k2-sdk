package stores

import (
	runtimeBases "github.com/tuxounet/k2-sdk/bases"

	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = storesTypes.ServiceKey

type Service struct {
	runtimeBases.BaseKernelService
}

func NewService(k runtimeTypes.IKernel) runtimeTypes.IKernelService {

	base := runtimeBases.NewBaseKernelService(k, ServiceKey)
	return &Service{base}

}

func (c *Service) GetStores() ([]string, error) {

	stores := c.getStores()

	storeNames := []string{}
	for _, store := range stores {
		storeNames = append(storeNames, store.Name)
	}

	return storeNames, nil

}

func (c *Service) GetStore(name string) (*storesTypes.Store, error) {

	stores := c.getStores()

	var found *storesTypes.Store
	for _, store := range stores {
		if store.Name == name {
			found = store
			break
		}
	}
	if found != nil {
		backends, err := c.GetBackends()
		if err != nil {
			c.GetLogger().ErrorF("Failed to get backends %s", err.Error())
			return nil, err
		}

		err = found.ResolveBackend(backends)
		if err != nil {
			c.GetLogger().ErrorF("Failed to resolve backend %s", err.Error())
			return nil, err
		}
		return found, nil
	}

	return nil, nil
}

func (c *Service) UpsertStore(store *storesTypes.Store) error {

	stores := c.getStores()

	wrote := false
	for i, s := range stores {
		if s.Name == store.Name {
			stores[i] = store
			wrote = true
			break
		}
	}
	if !wrote {
		stores = append(stores, store)
	}

	c.setStores(stores)

	return nil

}
