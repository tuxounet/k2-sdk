package k

import (
	"github.com/tuxounet/k2-sdk/kernel"
	"github.com/tuxounet/k2-sdk/types"
)

func HostApp(app types.IApp) {

	kernelRuntime := kernel.NewKernelRuntime(app)
	err := kernelRuntime.Init()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel init failed: %s", err.Error())
	}

	err = kernelRuntime.Register()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel register failed: %s", err.Error())
	}

	err = kernelRuntime.Start()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel start failed: %s", err.Error())
	}
	err = kernelRuntime.ListenAndServe()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel listenAndServe failed: %s", err.Error())
	}

}
