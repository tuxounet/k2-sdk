package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *BaseObjectStore[R]) getStore() (*types.Store, error) {

	storesServices := s.kernel.GetService(types.ServiceKey).(types.IStoreService)
	return storesServices.GetStore(s.storeName)

}
