package app

import (
	"embed"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"	
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
	"github.com/tuxounet/k2-sdk/samples/standalone/ui"
	{{ for k, v := range .components }}	
	{{ .k }} "{{ .module }}/components/{{ .k }}"
	{{ end }}
)

//go:embed config/*.yaml
var conf embed.FS

func NewApp() runtimeTypes.IApp {
	return runtimeBases.NewBaseApp(
		AppName,
		AppVersion,
		nil,
		&ui.Dist,
		&conf,
		[]runtimeTypes.AppComponentCtor{
			{{ for k, v := range .components }}
			{{ .k }}.NewComponent,
			{{ end }}
		},
	)
}
