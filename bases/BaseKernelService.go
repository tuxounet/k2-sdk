package bases

import (
	"github.com/tuxounet/k2-sdk/types"
)

type BaseKernelService struct {
	name   string
	kernel types.IKernel
	log    types.ILogger
	config map[string]string
	data   map[string]interface{}
}

func NewBaseKernelService(kernel types.IKernel, name string) BaseKernelService {
	log := kernel.GetLogger().CreateSubLogger(name)
	instance := BaseKernelService{
		name:   name,
		kernel: kernel,
		log:    log,
		config: make(map[string]string),
		data:   make(map[string]interface{}),
	}
	log.DebugF("Service %s created", name)
	return instance
}

func (b *BaseKernelService) GetName() string {
	return b.name
}

func (b *BaseKernelService) GetKernel() types.IKernel {
	return b.kernel
}

func (b *BaseKernelService) GetLogger() types.ILogger {
	return b.log
}

func (b *BaseKernelService) GetConfig(key string) string {
	return b.config[key]
}

func (b *BaseKernelService) SetConfig(key string, value string) {
	b.config[key] = value
}

func (b *BaseKernelService) GetData(key string) interface{} {
	return b.data[key]
}

func (b *BaseKernelService) SetData(key string, value interface{}) {
	b.data[key] = value
}

func (b *BaseKernelService) Init() error {
	return nil
}

func (b *BaseKernelService) Register() error {
	return nil
}

func (b *BaseKernelService) Start() error {
	return nil
}

func (b *BaseKernelService) Stop() error {
	return nil
}
