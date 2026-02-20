//go:build plugin

package main

import (
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"

	"github.com/tuxounet/k2-sdk/samples/with_plugins/plugins/plugin_one/temps"
	"github.com/tuxounet/k2-sdk/samples/with_plugins/plugins/plugin_one/ui"
)

// NewComponent is the exported symbol that the plugin host will look up
// This function creates and returns a new component instance
func NewComponent(app types.IApp) types.IAppComponent {
	return bases.NewBaseAppComponent(
		app,
		"plugin_one",
		1000,
		nil,
		nil,
		&ui.Dist,
		types.AccessPolicyPublic,
		[]types.AppControllerCtor{
			temps.NewController,
		},
	)
}
