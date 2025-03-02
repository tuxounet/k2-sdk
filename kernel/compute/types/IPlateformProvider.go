package types

import (
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type IBasePlateformProvider interface {
	GetName() string
	GetService() runtimeTypes.IKernelService
	GetLogger() runtimeTypes.ILogger
	Nuke() error
	Setup() error
	Render() ([]RunnerDefinition, error)
}

type IPlateformProvider[D any] interface {
	IBasePlateformProvider

	GetDefinitions() []D
	RegisterDefinition(definition D) error
}
