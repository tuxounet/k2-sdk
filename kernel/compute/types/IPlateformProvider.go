package types

import (
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type IBasePlateformProvider interface {
	GetName() string
	GetService() runtimeTypes.IKernelService
	GetLogger() runtimeTypes.ILogger
	Init() error
	Nuke() error
	Setup() error
	Render() ([]RunnerDefinition, error)
	Start() error
	Stop() error
}

type IPlateformProvider[D any] interface {
	IBasePlateformProvider

	GetDefinitions() []D
	RegisterDefinition(definition D) error
}
