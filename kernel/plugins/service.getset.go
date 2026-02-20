package plugins

import (
	"github.com/tuxounet/k2-sdk/types"
)

// GetLoadedPlugins returns the list of loaded plugins
func (s *PluginsService) GetLoadedPlugins() []types.PluginInfo {
	data := s.GetData("loadedPlugins")
	if data == nil {
		return []types.PluginInfo{}
	}
	return data.([]types.PluginInfo)
}
