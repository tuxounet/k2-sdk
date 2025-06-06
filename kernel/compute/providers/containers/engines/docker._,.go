package engines

import (
	"fmt"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type DockerEngine struct {
	Name            string
	logger          runtimeTypes.ILogger
	service         runtimeTypes.IKernelService
	portMapStore    storesTypes.IBaseObjectStore[[]types.PortsMapRecord]
	ingressRegistar ingressTypes.IngressRegistarFunction
}

func NewDockerEngine(service runtimeTypes.IKernelService, portMapStore storesTypes.IBaseObjectStore[[]types.PortsMapRecord], ingressRegistar ingressTypes.IngressRegistarFunction) types.IContainerEngine {
	log := service.GetLogger().CreateSubLogger(fmt.Sprintf("engine-%s", "docker"))

	return &DockerEngine{
		Name:            "docker",
		logger:          log,
		service:         service,
		ingressRegistar: ingressRegistar,
		portMapStore:    portMapStore,
	}
}

func (e *DockerEngine) Setup() error {
	e.logger.TraceF("Nuking up %s provider", e.Name)

	return nil
}
func (e *DockerEngine) Nuke() error {
	e.logger.TraceF("Setting up %s provider", e.Name)

	_, err := e.listContainers()
	if err != nil {
		e.logger.ErrorF("%s provider is not ready : %s", e.Name, err.Error())
		e.logger.InfoF("Try to install %s with the following command: 'sudo apt install %s'", e.Name, e.Name)
		return err
	}

	e.logger.DebugF("%s provider setup done", e.Name)
	return nil

}
