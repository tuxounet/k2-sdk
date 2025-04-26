package dbadmin

import (
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	bases.BaseAppController
}

func NewController(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "dbadmin", 22, nil)
	return &Controller{
		base,
	}
}
