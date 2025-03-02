package bases

import (
	"fmt"

	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type BaseObjectStore[R any] struct {
	storeKey     string
	storeName    string
	valueKey     string
	defaultValue string
	log          runtimeTypes.ILogger
	kernel       runtimeTypes.IKernel
}

func NewObjectStore[R any](kernel runtimeTypes.IKernel, parent runtimeTypes.ILoggable, storeName string, valueKey string, defaultValue string) storesTypes.IBaseObjectStore[R] {
	storeKey := fmt.Sprintf("%s.%s", storeName, valueKey)
	log := parent.GetLogger().CreateSubLogger(storeKey)
	return &BaseObjectStore[R]{
		storeName:    storeName,
		valueKey:     valueKey,
		storeKey:     storeKey,
		defaultValue: defaultValue,
		log:          log,
		kernel:       kernel,
	}
}
