package compute

import "github.com/tuxounet/k2-sdk/system"

func (s *Service) renderInventory() error {

	paths := s.getPathsService()
	inventoryPath := paths.CominePath("etc", "compute", "inventory")
	store, err := s.getRootStore()
	if err != nil {
		s.GetLogger().ErrorF("Failed to get root store: %s", err)
		return err
	}

	hosts := "all:\n"
	hosts += "  vars:\n"
	hosts += "    kernel:\n"
	hosts += "      name: " + s.GetKernel().GetApp().GetName() + "\n"
	hosts += "      version: " + s.GetKernel().GetApp().GetVersion() + "\n"
	err = store.WriteObject(paths.CominePath(inventoryPath, "hosts.yaml"), []byte(hosts))
	if err != nil {
		s.GetLogger().ErrorF("Failed to write hosts.yaml: %s", err)
		return err
	}

	appsConfigMap := s.getConfigService().GetCurrent()

	serializedConfig, err := system.DumpToYamlString(appsConfigMap)
	if err != nil {
		s.GetLogger().ErrorF("Failed to serialize full config map: %s", err)
		return err
	}

	err = store.WriteObject(paths.CominePath(inventoryPath, "group_vars", "all.yaml"), []byte(serializedConfig))
	if err != nil {
		s.GetLogger().ErrorF("Failed to write group_vars/all.yaml: %s", err)
		return err
	}

	return nil
}

func (s *Service) nukeInventory() error {
	paths := s.getPathsService()
	inventoryPath := paths.CominePath("etc", "compute", "inventory")
	store, err := s.getRootStore()
	if err != nil {
		s.GetLogger().ErrorF("Failed to get root store: %s", err.Error())
		return err
	}

	err = store.DeleteObject(inventoryPath)
	if err != nil {
		s.GetLogger().ErrorF("Failed to delete inventory: %s", err.Error())
		return err
	}

	return nil
}
