package app

 
{{ $module := .module }}

import (
	"embed"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"	
	runtimeTypes "github.com/tuxounet/k2-sdk/types"	 
	component "{{ $module }}/components/component"
)

//go:embed config/*.yaml
var conf embed.FS

func NewApp() runtimeTypes.IApp {
	return runtimeBases.NewBaseApp(
		AppName,
		AppVersion,
		nil,
		nil,
		&conf,
		[]runtimeTypes.AppComponentCtor{
			component.NewComponent,
		},
	)
}
