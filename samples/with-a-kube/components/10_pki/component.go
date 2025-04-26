package pki

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"

	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/10_pki/tls_ca"
	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/10_pki/tls_ingress"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed *.yaml
var conf embed.FS

func NewComponent(app types.IApp) types.IAppComponent {
	return bases.NewBaseAppComponent(
		app,
		"pki",
		10,
		nil,
		nil,
		&conf,
		types.AccessPolicyAdmin,
		[]types.AppControllerCtor{
			tls_ca.NewController,
			tls_ingress.NewController,
		},
	)
}
