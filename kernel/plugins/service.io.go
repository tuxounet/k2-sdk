package plugins

import "embed"

func (s *PluginsService) getAppExternals() *embed.FS {
	return s.GetKernel().GetApp().GetExternals()
}
