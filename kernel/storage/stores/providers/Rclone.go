package providers

import (
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/bases"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/types"

	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type RClone struct {
	bases.BaseStoreProvider
}

func NewRClone(service runtimeTypes.IKernelService) *RClone {
	base := bases.NewBaseStoreProvider(service, "rclone")

	return &RClone{base}
}

func (l *RClone) Setup() error {

	return nil
}

func (l *RClone) Exists(store *types.Store, path string) (bool, error) {

	return false, nil
}

func (l *RClone) Read(store *types.Store, path string) ([]byte, error) {

	return []byte{}, nil
}

func (l *RClone) Write(store *types.Store, path string, data []byte) error {

	return nil
}

func (l *RClone) Delete(store *types.Store, path string) error {

	return nil
}
