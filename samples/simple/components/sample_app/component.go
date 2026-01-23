package sample_app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/samples/simple/components/sample_app/controllers/files"
	"github.com/tuxounet/k2-sdk/samples/simple/components/sample_app/ui"

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

			files.NewFilesController,
		},
	)
}
