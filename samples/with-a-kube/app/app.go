package app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"

	provision "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/10_provision"
	auth "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/15_auth"
	sample_kube_app "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/20_sample_kube_app"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed config/*.yaml
var conf embed.FS

func NewApp() types.IApp {
	return bases.NewBaseApp(
		AppName,
		AppVersion,
		nil,
		nil,
		&conf,
		[]types.AppComponentCtor{
			provision.NewComponent,

			auth.NewComponent,
			sample_kube_app.NewComponent,
		},
	)
}
