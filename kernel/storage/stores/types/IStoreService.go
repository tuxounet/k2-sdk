package types

import runtimeTypes "github.com/tuxounet/k2-sdk/types"

const ServiceKey = "storage.stores"

type IStoreService interface {
	runtimeTypes.IKernelService
	GetStore(name string) (*Store, error)
}
