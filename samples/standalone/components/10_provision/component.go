package provision

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/10_provision/ansible"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/10_provision/browser"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/10_provision/tls_ca"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/10_provision/tls_ingress"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/10_provision/tty"

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
		types.AccessPolicyAuthenticated,
		[]types.AppControllerCtor{
			tls_ca.NewController,
			tls_ingress.NewController,
			browser.NewController,
			ansible.NewController,
			tty.NewController,
		},
	)
}
