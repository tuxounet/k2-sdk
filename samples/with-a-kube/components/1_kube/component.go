package kube

import (
	runtimeBases "github.com/tuxounet/k2-sdk/bases"

	tls_ca_issuer "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/10_tls_ca_issuer"
	ingress "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/1_ingress"
	cert_manager "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/2_cert-manager"
	tls_ca "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/3_tls_ca"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

func NewComponent(app runtimeTypes.IApp) runtimeTypes.IAppComponent {
	return runtimeBases.NewBaseAppComponent(
		app,
		"kube",
		1,
		nil,
		nil,
		nil,
		[]runtimeTypes.AppControllerCtor{
			ingress.NewController,
			cert_manager.NewController,
			tls_ca.NewController,
			tls_ca_issuer.NewController,
		},
	)
}
