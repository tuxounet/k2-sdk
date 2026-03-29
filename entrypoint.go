package k

import (
	_ "embed"
	"os"

	"github.com/tuxounet/k2-sdk/kernel"
	"github.com/tuxounet/k2-sdk/kernel/compute"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	"github.com/tuxounet/k2-sdk/types"
)

//go:embed version.txt
var RuntimeVersion string

func HostApp(app types.IApp) {

	kernelRuntime := kernel.NewKernelRuntime(app, RuntimeVersion, false)
	applyForceComputeFlag(kernelRuntime)
	doRun(kernelRuntime)

}
func HostUnsecureApp(app types.IApp) {
	kernelRuntime := kernel.NewKernelRuntime(app, RuntimeVersion, true)
	applyForceComputeFlag(kernelRuntime)
	doRun(kernelRuntime)

}

func HostProvisionOnly(app types.IApp) {
	kernelRuntime := kernel.NewKernelRuntime(app, RuntimeVersion, false)
	kernelRuntime.SetForceCompute(true)

	err := kernelRuntime.Init()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel init failed: %s", err.Error())
	}
	err = kernelRuntime.Register()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel register failed: %s", err.Error())
	}

	doProvisionOnly(kernelRuntime)
}

func HostTeardownOnly(app types.IApp) {
	kernelRuntime := kernel.NewKernelRuntime(app, RuntimeVersion, false)
	kernelRuntime.SetForceCompute(true)

	err := kernelRuntime.Init()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel init failed: %s", err.Error())
	}
	err = kernelRuntime.Register()
	if err != nil {
		kernelRuntime.GetLogger().PanicF("kernel register failed: %s", err.Error())
	}

	doTeardownOnly(kernelRuntime)
}

func applyForceComputeFlag(rt *kernel.KernelRuntime) {
	for _, arg := range os.Args[1:] {
		if arg == "--force-compute" {
			rt.SetForceCompute(true)
			return
		}
	}
}

func hasFlag(flag string) bool {
	for _, arg := range os.Args[1:] {
		if arg == flag {
			return true
		}
	}
	return false
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

	if hasFlag("--provision-only") {
		doProvisionOnly(rt)
		return
	}

	if hasFlag("--teardown-only") {
		doTeardownOnly(rt)
		return
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

func doProvisionOnly(rt *kernel.KernelRuntime) {
	computeService := rt.GetService(types.KernelServiceContextKey(compute.ServiceKey)).(*compute.Service)

	err := computeService.ExecVerb(computeTypes.RunnerVerbProvision)
	if err != nil {
		rt.GetLogger().PanicF("provision failed: %s", err.Error())
	}

	rt.GetLogger().InfoF("provision completed successfully")
}

func doTeardownOnly(rt *kernel.KernelRuntime) {
	computeService := rt.GetService(types.KernelServiceContextKey(compute.ServiceKey)).(*compute.Service)

	err := computeService.ExecVerb(computeTypes.RunnerVerbStop)
	if err != nil {
		rt.GetLogger().WarnF("stop phase failed (continuing to teardown): %s", err.Error())
	}

	err = computeService.ExecVerb(computeTypes.RunnerVerbTeardown)
	if err != nil {
		rt.GetLogger().PanicF("teardown failed: %s", err.Error())
	}

	rt.GetLogger().InfoF("teardown completed successfully")
}
