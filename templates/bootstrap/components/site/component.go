package site

import (
	"{{ .module }}/site/health"
	runtimeBases "{{ .sdk_module }}/bases"
	runtimeTypes "{{ .sdk_module }}/types"
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
		[]runtimeTypes.AppControllerCtor{
			health.NewController,
		},
	)
}
