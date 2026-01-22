package kube

import (
	runtimeBases "github.com/tuxounet/k2-sdk/bases"

	// tls_ca "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/10_tls_ca"
	// tls_ca_issuer "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/11_tls_ca_issuer"
	// cert_manager "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/12_cert-manager"
	loadbalancer "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/1_loadbalancer"
	reverse_proxy "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/2_reverse_proxy"
	ingress "github.com/tuxounet/k2-sdk/samples/with-a-kube/components/1_kube/5_ingress"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

func NewComponent(app runtimeTypes.IApp) runtimeTypes.IAppComponent {
	return runtimeBases.NewBaseAppComponent(
		app,
		"kube",
		50,
		nil,
		nil,
		nil,
		runtimeTypes.AccessPolicyPublic,
		[]runtimeTypes.AppControllerCtor{
			loadbalancer.NewController,
			reverse_proxy.NewController,

			ingress.NewController,
			// cert_manager.NewController,
			// tls_ca.NewController,
			// tls_ca_issuer.NewController,
		},
	)
}
