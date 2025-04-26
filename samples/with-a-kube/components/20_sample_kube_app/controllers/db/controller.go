package db

import (
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	bases.BaseAppController
}

func NewController(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "db", 21, nil)
	return &Controller{
		base,
	}
}
