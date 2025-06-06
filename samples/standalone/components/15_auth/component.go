package auth

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/15_auth/controllers/auth"

	"github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app/ui"

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
		&ui.Dist,
		&conf,

		[]types.AppControllerCtor{
			auth.NewController,
		},
	)
}
