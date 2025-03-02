package site

import (
	"craft.f.krux/git/runtimes/k-app-runtime.git/types"
	runtimeBases "github.com/tuxounet/k2-sdk/bases"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
	"{{ .module }}/site/health"
)

func NewComponent(app runtimeTypes.IApp) runtimeTypes.IAppComponent {
	return runtimeBases.NewBaseAppComponent(
		app,
		"site",
		1,
		nil,
		nil,
		nil,
		runtimeTypes.AccessPolicyPublic,
		[]types.AppControllerCtor{
			health.NewController,
		},
	)
}
