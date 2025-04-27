package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	runtimesTypes "github.com/tuxounet/k2-sdk/types"
)

type BasePlateformProvider[D any] struct {
	name            string
	service         runtimesTypes.IKernelService
	log             runtimesTypes.ILogger
	definitions     []D
	runners         []types.RunnerDefinition
	data            map[string]interface{}
	ingressRegistar ingressTypes.IngressRegistarFunction
	isRequired      bool
}

func NewBasePlateformProvider[D any](service runtimesTypes.IKernelService, key string, ingressRegistar ingressTypes.IngressRegistarFunction) BasePlateformProvider[D] {
	log := service.GetLogger().CreateSubLogger(key)
	instance := BasePlateformProvider[D]{
		name:            key,
		service:         service,
		log:             log,
		isRequired:      false,
		data:            make(map[string]interface{}),
		ingressRegistar: ingressRegistar,
	}

	instance.definitions = make([]D, 0)
	instance.runners = make([]types.RunnerDefinition, 0)
	return instance
}
