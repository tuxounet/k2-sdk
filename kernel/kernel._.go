package kernel

import (
	"context"
	"os"
	"path"

	"github.com/tuxounet/k2-sdk/kernel/app"
	"github.com/tuxounet/k2-sdk/kernel/compute"
	"github.com/tuxounet/k2-sdk/kernel/monitoring/logging"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress"
	"github.com/tuxounet/k2-sdk/kernel/plugins"
	"github.com/tuxounet/k2-sdk/kernel/scheduler"

	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/profile"
	"github.com/tuxounet/k2-sdk/kernel/secrets"
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores"
	"github.com/tuxounet/k2-sdk/kernel/storage/volumes"
	"github.com/tuxounet/k2-sdk/types"
)

type KernelRuntime struct {
	hostName    string
	hostVersion string
	rootContext context.Context
	app         types.IApp
	runDir      string
	log         types.ILogger
	services    []types.KernelServiceContextKey
	unsecure    bool
}

func NewKernelRuntime(hostedApp types.IApp, hostVersion string, unsecure bool) *KernelRuntime {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootContext := context.Background()

	runDir := path.Join(cwd, ".data")
	if os.Getenv("RUN_DIR") != "" {
		runDir = os.Getenv("RUN_DIR")
	}
	err = os.MkdirAll(runDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	kernel := &KernelRuntime{
		hostName:    hostedApp.GetName(),
		hostVersion: hostVersion,
		rootContext: rootContext,
		app:         hostedApp,
		runDir:      runDir,
		unsecure:    unsecure,
	}

	//Early logging service init
	loggingService := logging.NewService(kernel)
	kernel.SetService(loggingService)

	//early init
	err = loggingService.Init()
	if err != nil {
		panic(err)
	}
	kernel.log = loggingService.GetRootLogger()
	kernel.log.TraceF("[BOOT] %s %s created", kernel.hostName, kernel.hostVersion)

	servicesCreateList := [](func(types.IKernel) types.IKernelService){
		paths.NewService,
		profile.NewService,
		volumes.NewService,
		stores.NewService,
		config.NewService,
		secrets.NewService,
		compute.NewService,
		plugins.NewService,
		app.NewService,
		scheduler.NewService,
		ingress.NewService,
	}
	serviceList := make([]types.KernelServiceContextKey, 0)

	for _, call := range servicesCreateList {
		service := call(kernel)
		if service != nil {
			kernel.SetService(service)
			serviceList = append(serviceList, types.KernelServiceContextKey(service.GetName()))
		} else {
			kernel.GetLogger().PanicF("service not created during %s", call)

		}
	}
	kernel.services = serviceList
	kernel.log.TraceF("[BOOT] complete")
	return kernel
}
