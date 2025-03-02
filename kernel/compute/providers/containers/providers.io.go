package containers

import (
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores"
	storeTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (p *Provider) getPathsService() *paths.Service {
	return p.GetService().GetKernel().GetService(paths.ServiceKey).(*paths.Service)
}

func (p *Provider) getRootStore() (*storeTypes.Store, error) {
	stores := p.GetService().GetKernel().GetService(stores.ServiceKey).(*stores.Service)
	store, err := stores.GetStore("root")
	if err != nil {
		return nil, err
	}
	return store, nil
}
