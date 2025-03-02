package sample_app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app/controllers/files"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app/controllers/httpbin"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app/controllers/sample"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app/controllers/sample2"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app/ui"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed *.yaml
var conf embed.FS

func NewComponent(app types.IApp) types.IAppComponent {
	return bases.NewBaseAppComponent(
		app,
		"sample_app",
		20,
		nil,
		&ui.Dist,
		&conf,
		types.AccessPolicyAuthenticated,
		[]types.AppControllerCtor{
			sample.NewSampleController,
			sample2.NewSample2Controller,
			files.NewFilesController,
			httpbin.NewController,
		},
	)
}
