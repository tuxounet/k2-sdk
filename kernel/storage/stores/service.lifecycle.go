package stores

import (
	"fmt"
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
			s.GetLogger().ErrorF("failed to setup storage backend %s: %s", backend.GetName(), err.Error())
			return fmt.Errorf("storage: failed to setup backend %s: %w", backend.GetName(), err)
		}
	}

	runDir := s.GetKernel().GetRunDirectory()
	profiles := s.getProfileService()
	userDir, err := profiles.GetUserDirectory()
	if err != nil {
		s.GetLogger().ErrorF("failed to get user directory: %s", err.Error())
		return fmt.Errorf("storage: failed to get user directory: %w", err)
	}

	err = s.UpsertStore(&types.Store{
		Name:    "root",
		Backend: "local",
		Flags: map[string]any{
			"path": runDir,
		},
	})
	if err != nil {
		s.GetLogger().ErrorF("failed to upsert root store: %s", err.Error())
		return fmt.Errorf("storage: failed to upsert root store: %w", err)
	}

	err = s.UpsertStore(&types.Store{
		Name:    "local",
		Backend: "local",
		Flags: map[string]any{
			"path": userDir,
		},
	})
	if err != nil {
		s.GetLogger().ErrorF("failed to upsert local store: %s", err.Error())
		return fmt.Errorf("storage: failed to upsert local store: %w", err)
	}

	rootStore, err := s.GetStore("root")
	if err != nil {
		s.GetLogger().ErrorF("failed to get root store: %s", err.Error())
		return fmt.Errorf("storage: failed to get root store: %w", err)
	}

	//cleanup temp dir
	err = rootStore.DeleteObject("tmp")
	if err != nil {
		s.GetLogger().ErrorF("failed to delete tmp directory: %s", err.Error())
		return fmt.Errorf("storage: failed to cleanup tmp directory: %w", err)
	}
	err = rootStore.WriteObject(filepath.Join("tmp", ".keep"), []byte{})
	if err != nil {
		s.GetLogger().ErrorF("failed to create tmp directory: %s", err.Error())
		return fmt.Errorf("storage: failed to create tmp directory: %w", err)
	}
	return nil
}
