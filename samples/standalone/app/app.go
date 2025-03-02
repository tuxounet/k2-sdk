package app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	provision "github.com/tuxounet/k2-sdk/samples/standalone/components/10_provision"
	auth "github.com/tuxounet/k2-sdk/samples/standalone/components/15_auth"
	sample_app "github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app"
	"github.com/tuxounet/k2-sdk/samples/standalone/ui"

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
			provision.NewComponent,
			auth.NewComponent,
			sample_app.NewComponent,
		},
	)
}
