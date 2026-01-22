package auth

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	auth "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/15_auth/controllers/auth"

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
			auth.NewController,
		},
	)
}
