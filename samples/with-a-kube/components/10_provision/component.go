package provision

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"

	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/10_provision/ca_issuer"
	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/10_provision/tls_ca"
	"github.com/tuxounet/k2-sdk/samples/with-a-kube/components/10_provision/tls_ingress"

	"github.com/tuxounet/k2-sdk/types"
)

//go:embed *.yaml
var conf embed.FS

func NewComponent(app types.IApp) types.IAppComponent {
	return bases.NewBaseAppComponent(
		app,
		"provision",
		10,
		nil,
		nil,
		&conf,
		[]types.AppControllerCtor{
			tls_ca.NewController,
			tls_ingress.NewController,
			ca_issuer.NewController,
		},
	)
}
