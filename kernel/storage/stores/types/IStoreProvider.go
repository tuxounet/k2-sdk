package types

import "github.com/tuxounet/k2-sdk/types"

type IStoreProvider interface {
	GetName() string
	Setup() error
	GetLogger() types.ILogger
	Exists(store *Store, path string) (bool, error)
	Read(store *Store, path string) ([]byte, error)
	Write(store *Store, path string, data []byte) error
	Delete(store *Store, path string) error
}
