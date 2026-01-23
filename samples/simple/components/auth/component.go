package auth

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"

	"github.com/tuxounet/k2-sdk/samples/simple/components/auth/controllers/users"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed *.yaml
var conf embed.FS

func NewComponent(app types.IApp) types.IAppComponent {
	return bases.NewBaseAppComponent(
		app,
		"auth",
		15,
		nil,
		nil,
		&conf,
		types.AccessPolicyPublic,
		[]types.AppControllerCtor{
			users.NewController,
		},
	)
}
