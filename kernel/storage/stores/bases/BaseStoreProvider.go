package bases

import (
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type BaseStoreProvider struct {
	name    string
	service runtimeTypes.IKernelService
	log     runtimeTypes.ILogger
}

func NewBaseStoreProvider(service runtimeTypes.IKernelService, name string) BaseStoreProvider {
	log := service.GetLogger().CreateSubLogger(name)

	base := BaseStoreProvider{
		name:    name,
		log:     log,
		service: service,
	}
	log.DebugF("created new backend %s", name)

	return base
}

func (b BaseStoreProvider) GetName() string {
	return b.name
}

func (b BaseStoreProvider) GetLogger() runtimeTypes.ILogger {
	return b.log
}

func (b BaseStoreProvider) Setup() error {
	return nil
}
