package kernel

import (
	"context"

	"github.com/tuxounet/k2-sdk/types"
)

func (g *KernelRuntime) GetName() string {
	return g.hostName
}

func (g *KernelRuntime) GetVersion() string {
	return g.hostVersion
}

func (k *KernelRuntime) GetRootContext() context.Context {
	return k.rootContext
}

func (g *KernelRuntime) GetApp() types.IApp {
	return g.app
}

func (g *KernelRuntime) GetRunDirectory() string {
	return g.runDir
}

func (g *KernelRuntime) GetLogger() types.ILogger {
	return g.log
}

func (g *KernelRuntime) GetService(key types.KernelServiceContextKey) types.IKernelService {
	service := g.rootContext.Value(key)
	if service == nil {
		g.log.WarnF("Service not found for key %s", key)
		return nil
	}

	return service.(types.IKernelService)
}

func (g *KernelRuntime) SetService(service types.IKernelService) {
	key := types.KernelServiceContextKey(service.GetName())
	g.rootContext = context.WithValue(g.rootContext, key, service)
}

func (g *KernelRuntime) IsUnsecure() bool {
	return g.unsecure
}
