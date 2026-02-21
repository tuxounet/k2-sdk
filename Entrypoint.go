package k

import (
	_ "embed"

	"github.com/tuxounet/k2-sdk/kernel"
	"github.com/tuxounet/k2-sdk/types"
)

//go:embed version.txt
var RuntimeVersion string

func HostApp(app types.IApp) {

	kernelRuntime := kernel.NewKernelRuntime(app, RuntimeVersion, false)
	doRun(kernelRuntime)

}
func HostUnsecureApp(app types.IApp) {
	kernelRuntime := kernel.NewKernelRuntime(app, RuntimeVersion, true)
	doRun(kernelRuntime)

}
func doRun(rt *kernel.KernelRuntime) {
	err := rt.Init()
	if err != nil {
		rt.GetLogger().PanicF("kernel init failed: %s", err.Error())
	}

	err = rt.Register()
	if err != nil {
		rt.GetLogger().PanicF("kernel register failed: %s", err.Error())
	}

	err = rt.Start()
	if err != nil {
		rt.GetLogger().PanicF("kernel start failed: %s", err.Error())
	}
	err = rt.ListenAndServe()
	if err != nil {
		rt.GetLogger().PanicF("kernel listenAndServe failed: %s", err.Error())
	}

}
