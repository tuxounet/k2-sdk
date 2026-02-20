package stores

import (
	"path/filepath"

	"github.com/tuxounet/k2-sdk/kernel/storage/stores/providers"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Service) Init() error {

	backends := []types.IStoreProvider{
		providers.NewLocal(s),
		providers.NewRClone(s),
	}
	s.setBackends(backends)

	for _, backend := range backends {
		err := backend.Setup()
		if err != nil {
			s.GetLogger().ErrorF("Failed to setup backend %s", err.Error())
			return err
		}
	}

	runDir := s.GetKernel().GetRunDirectory()
	profiles := s.getProfileService()
	userDir, err := profiles.GetUserDirectory()
	if err != nil {
		s.GetLogger().ErrorF("Failed to get user directory %s", err.Error())
		return err
	}

	err = s.UpsertStore(&types.Store{
		Name:    "root",
		Backend: "local",
		Flags: map[string]any{
			"path": runDir,
		},
	})
	if err != nil {
		s.GetLogger().ErrorF("Failed to upsert store %s", err.Error())
		return err
	}

	err = s.UpsertStore(&types.Store{
		Name:    "local",
		Backend: "local",
		Flags: map[string]any{
			"path": userDir,
		},
	})
	if err != nil {
		s.GetLogger().ErrorF("Failed to upsert store %s", err.Error())
		return err
	}

	rootStore, err := s.GetStore("root")
	if err != nil {
		s.GetLogger().ErrorF("Failed to get root store %s", err.Error())
		return err
	}

	//cleanup temp dir
	err = rootStore.DeleteObject("tmp")
	if err != nil {
		s.GetLogger().ErrorF("Failed to delete tmp directory %s", err.Error())
		return err
	}
	err = rootStore.WriteObject(filepath.Join("tmp", ".keep"), []byte{})
	if err != nil {
		s.GetLogger().ErrorF("Failed to create tmp directory %s", err.Error())
		return err
	}
	return nil
}
