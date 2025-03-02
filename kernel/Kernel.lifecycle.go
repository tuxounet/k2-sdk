package kernel

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tuxounet/k2-sdk/kernel/network/ingress"

	"github.com/tuxounet/k2-sdk/types"
)

func (k *KernelRuntime) Init() error {
	k.GetLogger().TraceF("begin init")

	for _, service := range k.services {
		k.GetLogger().TraceF("init service %s", service)
		err := k.GetService(service).Init()
		if err != nil {
			k.GetLogger().ErrorF("service %s init failed: %s", service, err.Error())
			return err
		}
		k.GetLogger().TraceF("inited service %s", service)
	}

	k.GetLogger().TraceF("end init")
	return nil
}

func (k *KernelRuntime) Register() error {
	k.GetLogger().TraceF("begin register")

	for _, key := range k.services {
		k.GetLogger().TraceF("register service %s", key)
		service := k.GetService(key)
		err := service.Register()
		if err != nil {
			k.GetLogger().ErrorF("service %s register failed: %s", key, err.Error())
			return err
		}
		k.GetLogger().TraceF("registered service %s", key)
	}

	k.GetLogger().TraceF("end register")
	return nil
}

func (k *KernelRuntime) Start() error {
	k.GetLogger().TraceF("begin start")

	for _, key := range k.services {
		k.GetLogger().TraceF("start service %s", key)
		service := k.GetService(key)
		err := service.Start()
		if err != nil {
			k.GetLogger().ErrorF("service %s register failed: %s", key, err.Error())
			return err
		}
		k.GetLogger().TraceF("started service %s", key)
	}

	k.GetLogger().TraceF("end start")
	return nil
}

func (k *KernelRuntime) Stop() error {
	k.GetLogger().TraceF("begin stop")

	for _, key := range k.services {
		k.GetLogger().TraceF("Stop service %s", key)
		service := k.GetService(key)
		err := service.Stop()
		if err != nil {
			k.GetLogger().ErrorF("service %s Stop failed: %s", key, err.Error())
			return err
		}
		k.GetLogger().TraceF("stopped service %s", key)
	}
	k.GetLogger().TraceF("end stop")
	return nil
}

func (k *KernelRuntime) ListenAndServe() error {
	k.GetLogger().TraceF("begin listen and serve")

	ingressHttpService := k.GetService(types.KernelServiceContextKey(ingress.ServiceKey)).(*ingress.Service)
	server := ingressHttpService.GetServer()
	if server == nil {
		k.GetLogger().Panic("ingress service server not found")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	k.log.InfoF("Kernel Running...")
	ingressHttpService.Listen()

	<-signalChan

	k.log.WarnF("Kill signal received, stopping...")
	err := k.Stop()
	if err != nil {
		log.Panicf("Kernel stop failed: %s", err.Error())
	}

	k.log.InfoF("Kernel stopped")
	return nil
}
