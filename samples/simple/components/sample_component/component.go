package sample_component

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/samples/simple/components/sample_component/controllers/files"
	"github.com/tuxounet/k2-sdk/samples/simple/components/sample_component/ui"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed *.yaml
var conf embed.FS

func NewComponent(app types.IApp) types.IAppComponent {
	return bases.NewBaseAppComponent(
		app,
		"sample_component",
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
