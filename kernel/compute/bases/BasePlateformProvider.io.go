package bases

import "github.com/tuxounet/k2-sdk/kernel/storage/paths"

func (s BasePlateformProvider[D]) getPaths() *paths.Service {
	return s.service.GetKernel().GetService(paths.ServiceKey).(*paths.Service)
}
