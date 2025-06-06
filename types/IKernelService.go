package types

type KernelServiceContextKey string

type IKernelService interface {
	ILoggable
	GetName() string
	GetKernel() IKernel
	GetConfig(key string) string
	SetConfig(key string, value string)
	GetData(key string) interface{}
	SetData(key string, value interface{})
	Init() error
	Register() error
	Start() error
	Stop() error
}
