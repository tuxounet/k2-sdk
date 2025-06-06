package frontend

import (
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	bases.BaseAppController
}

func NewController(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "frontend", 30, nil, types.AccessPolicyPublic)
	return &Controller{
		base,
	}
}
