package containers

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Provider) getPortsMapSore() storesTypes.IBaseObjectStore[[]types.PortsMapRecord] {
	return s.GetData("portsmap").(storesTypes.IBaseObjectStore[[]types.PortsMapRecord])

}

func (s *Provider) setContainerEngine(engine types.IContainerEngine) {
	s.SetData("engine", engine)
}
func (s *Provider) getContainerEngine() types.IContainerEngine {
	engine, ok := s.GetData("engine").(types.IContainerEngine)
	if !ok {
		return nil
	}
	return engine
}
