package auth_backend

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

//go:embed config/*.yaml
var conf embed.FS

type Controller struct {
	bases.BaseAppController
}

func NewController(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "auth_backend", 1, &conf)
	return &Controller{
		base,
	}
}
