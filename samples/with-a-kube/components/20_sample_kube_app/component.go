package sample_kube_app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"

	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/20_sample_kube_app/namespaces/admin"
	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/20_sample_kube_app/namespaces/database"
	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/20_sample_kube_app/namespaces/frontend"
	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/20_sample_kube_app/ui"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed *.yaml
var conf embed.FS

func NewComponent(app types.IApp) types.IAppComponent {
	return bases.NewBaseAppComponent(
		app,
		"sample_kube_app",
		20,
		nil,
		&ui.Dist,
		&conf,
		types.AccessPolicyPublic,
		[]types.AppControllerCtor{
			admin.NewController,
			database.NewController,
			frontend.NewController,
		},
	)
}
