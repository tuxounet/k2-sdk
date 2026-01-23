package app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"

	"github.com/tuxounet/k2-sdk/samples/simple/components/auth"
	"github.com/tuxounet/k2-sdk/samples/simple/components/sample_app"

	"github.com/tuxounet/k2-sdk/samples/simple/ui"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed config/*.yaml
var conf embed.FS

func NewApp() types.IApp {
	return bases.NewBaseApp(
		AppName,
		AppVersion,
		nil,
		&ui.Dist,
		&conf,
		[]types.AppComponentCtor{

			auth.NewComponent,
			sample_app.NewComponent,
		},
	)
}
