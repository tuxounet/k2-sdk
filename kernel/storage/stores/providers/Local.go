package providers

import (
	"os"
	"path"

	"github.com/tuxounet/k2-sdk/kernel/storage/stores/bases"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Local struct {
	bases.BaseStoreProvider
}

func NewLocal(service runtimeTypes.IKernelService) *Local {

	base := bases.NewBaseStoreProvider(service, "local")

	return &Local{base}

}

func (l *Local) getBasePath(store *types.Store) (string, error) {
	pathFlag := store.Flags["path"]
	if pathFlag == "" {
		l.GetLogger().ErrorF("store %s has no path flag", store.GetName())
		return "", nil
	}
	return pathFlag.(string), nil

}

func (l *Local) Setup() error {

	return nil
}

func (l *Local) Exists(store *types.Store, filePath string) (bool, error) {

	basePath, err := l.getBasePath(store)
	if err != nil {
		l.GetLogger().ErrorF("Failed to get base path %s", err.Error())
		return false, err
	}

	realFilePath := path.Join(basePath, filePath)

	_, err = os.Stat(realFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		l.GetLogger().ErrorF("Failed to stat file %s", err.Error())
		return false, err
	}

	return true, nil
}

func (l *Local) Read(store *types.Store, filePath string) ([]byte, error) {
	basePath, err := l.getBasePath(store)
	if err != nil {
		l.GetLogger().ErrorF("Failed to get base path %s", err.Error())
		return nil, err
	}

	realFilePath := path.Join(basePath, filePath)

	body, err := os.ReadFile(realFilePath)
	if err != nil {
		l.GetLogger().ErrorF("Failed to read file %s", err.Error())
		return nil, err
	}

	return body, nil
}

func (l *Local) Write(store *types.Store, filePath string, data []byte) error {
	basePath, err := l.getBasePath(store)
	if err != nil {
		l.GetLogger().ErrorF("Failed to get base path %s", err.Error())
		return err
	}

	realFilePath := path.Join(basePath, filePath)
	dirName := path.Dir(realFilePath)
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		l.GetLogger().ErrorF("Failed to create directory %s", err.Error())
		return err
	}

	err = os.WriteFile(realFilePath, data, 0644)
	if err != nil {
		l.GetLogger().ErrorF("Failed to write file %s", err.Error())
		return err
	}

	return nil
}

func (l *Local) Delete(store *types.Store, filePath string) error {
	basePath, err := l.getBasePath(store)
	if err != nil {
		l.GetLogger().ErrorF("Failed to get base path %s", err.Error())
		return err
	}

	realFilePath := path.Join(basePath, filePath)

	exists, err := l.Exists(store, filePath)
	if err != nil {
		l.GetLogger().ErrorF("Failed to check if file exists %s", err.Error())
		return err
	}
	if !exists {
		return nil
	}

	stat, err := os.Stat(realFilePath)
	if err != nil {
		l.GetLogger().ErrorF("Failed to stat file %s", err.Error())
		return err
	}
	if stat.IsDir() {
		err = os.RemoveAll(realFilePath)
		if err != nil {
			l.GetLogger().ErrorF("Failed to delete directory %s", err.Error())
			return err

		}
	} else {
		err = os.Remove(realFilePath)

		if err != nil {
			l.GetLogger().ErrorF("Failed to delete file %s", err.Error())
			return err
		}
	}

	return nil
}
