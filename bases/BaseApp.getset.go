package bases

import (
	"embed"

	"github.com/swaggo/swag"
	"github.com/tuxounet/k2-sdk/types"
)

func (a *BaseApp) GetLogger() types.ILogger {
	return a.log
}

func (a *BaseApp) SetLogger(logger types.ILogger) {
	a.log = logger
}

func (a *BaseApp) GetName() string {
	return a.name
}

func (a *BaseApp) GetVersion() string {
	return a.version
}

func (a *BaseApp) GetDocs() *swag.Spec {

	return a.docs
}

func (a *BaseApp) GetUI() *embed.FS {
	return a.ui
}

func (a *BaseApp) GetConfig() *embed.FS {
	return a.config
}

func (a *BaseApp) GetComponents() []types.IAppComponent {
	if a.components == nil {
		kernel := a.GetKernel()
		if kernel == nil {
			return nil
		}

		a.components = make([]types.IAppComponent, 0)
		for _, ctor := range a.componentsCtors {
			component := ctor(a)
			a.components = append(a.components, component)
		}
	}

	return a.components
}

func (a *BaseApp) GetComponent(name string) types.IAppComponent {
	for _, controller := range a.GetComponents() {
		if controller.GetName() == name {
			return controller
		}
	}

	return nil
}

func (a *BaseApp) GetKernel() types.IKernel {
	return a.kernel
}

func (a *BaseApp) SetKernel(kernel types.IKernel) {
	a.kernel = kernel
}
