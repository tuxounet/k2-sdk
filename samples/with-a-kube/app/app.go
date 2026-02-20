package app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"

	kube "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube"
	// auth "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/15_auth"
	// sample_kube_app "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/20_sample_kube_app"

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
			kube.NewComponent,
			// auth.NewComponent,
			// sample_kube_app.NewComponent,
		},
		nil,
	)
}
