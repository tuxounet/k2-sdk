package types

import (
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
)

type IContainerEngine interface {
	Setup() error
	Nuke() error
	RenderPlaybookTasks(definition ContainerDefinition, verb computeTypes.RunnerVerb) (string, error)
}
