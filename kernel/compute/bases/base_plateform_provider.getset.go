package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

func (b *BasePlateformProvider[D]) GetName() string {
	return b.name
}

func (b *BasePlateformProvider[D]) GetLogger() runtimeTypes.ILogger {
	return b.log
}

func (b *BasePlateformProvider[D]) GetService() runtimeTypes.IKernelService {
	return b.service
}

func (s *BasePlateformProvider[D]) GetIsRequired() bool {
	return s.isRequired
}
func (s *BasePlateformProvider[D]) SetIsRequired(isNeeded bool) {
	s.isRequired = isNeeded
}

func (b *BasePlateformProvider[D]) GetData(key string) interface{} {
	return b.data[key]
}

func (b *BasePlateformProvider[D]) SetData(key string, value interface{}) {
	b.data[key] = value
}

func (b *BasePlateformProvider[D]) GetDefinitions() []D {
	return b.definitions
}

func (b *BasePlateformProvider[D]) ResetDefinitions() {
	b.definitions = make([]D, 0)
}

func (b *BasePlateformProvider[D]) GetRunners() []types.RunnerDefinition {
	return b.runners
}

func (b *BasePlateformProvider[D]) GetIngressRegistar() ingressTypes.IngressRegistarFunction {
	return b.ingressRegistar
}
