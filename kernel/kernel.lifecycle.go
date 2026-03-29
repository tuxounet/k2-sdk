package kernel

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tuxounet/k2-sdk/kernel/network/ingress"

	"github.com/tuxounet/k2-sdk/types"
)

func (k *KernelRuntime) Init() error {
	k.GetLogger().TraceF("[INIT]")

	for _, service := range k.services {
		err := k.GetService(service).Init()
		if err != nil {
			k.GetLogger().ErrorF("service %s init failed: %s", service, err.Error())
			return err
		}
		k.GetLogger().TraceF("[INIT] inited service %s", service)
	}

	k.GetLogger().TraceF("[INIT] complete")
	return nil
}

func (k *KernelRuntime) Register() error {
	k.GetLogger().TraceF("[REGISTER]")

	for _, key := range k.services {
		service := k.GetService(key)
		err := service.Register()
		if err != nil {
			return fmt.Errorf("service %s register failed: %w", key, err)
		}
		k.GetLogger().TraceF("[REGISTER] registered service %s", key)
	}

	k.GetLogger().TraceF("[REGISTER] complete")
	return nil
}

func (k *KernelRuntime) Start() error {
	k.GetLogger().TraceF("[START]")

	for _, key := range k.services {
		service := k.GetService(key)
		err := service.Start()
		if err != nil {
			return fmt.Errorf("service %s start failed: %w", key, err)

		}
		k.GetLogger().TraceF("[START] started service %s", key)
	}

	k.GetLogger().TraceF("[START] complete")
	return nil
}

func (k *KernelRuntime) Stop() error {
	k.GetLogger().TraceF("[STOP]")

	for _, key := range k.services {
		service := k.GetService(key)
		err := service.Stop()
		if err != nil {
			k.GetLogger().ErrorF("service %s Stop failed: %s", key, err.Error())
			return err
		}
		k.GetLogger().TraceF("[STOP] stopped service %s", key)
	}
	k.GetLogger().TraceF("[STOP] complete")
	return nil
}

func (k *KernelRuntime) ListenAndServe() error {
	k.GetLogger().TraceF("[LISTEN] begin listen and serve")

	ingressHttpService := k.GetService(types.KernelServiceContextKey(ingress.ServiceKey)).(*ingress.Service)
	server := ingressHttpService.GetServer()
	if server == nil {
		return fmt.Errorf("ingress service not initialized")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	k.log.InfoF("[LISTEN] Kernel Running...")
	err := ingressHttpService.Listen()
	if err != nil {
		return fmt.Errorf("ingress service listen failed: %s", err.Error())
	}

	<-signalChan

	k.log.WarnF("[LISTEN] Kill signal received, stopping...")
	err = k.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop kernel: %w", err)
	}

	k.log.InfoF("[LISTEN] Kernel stopped")
	return nil
}
