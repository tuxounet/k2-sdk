{{ $componentList := split "," .components }}
{{ $module := .module }}
package app


import (
	"embed"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"	
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
	"github.com/tuxounet/k2-sdk/samples/standalone/ui"
	{{range $componentList := .}}
	{{ . }} "{{ $module }}/components/{{ . }}"
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
			{{range $componentList := .}}
			{{ . }}.NewComponent,
			{{ end }}
		},
	)
}
