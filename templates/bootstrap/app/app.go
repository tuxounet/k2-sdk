{{ $componentList := splitList  "," .components }}
{{ $module := .module }}
package app


import (
	"embed"

	runtimeBases "github.com/tuxounet/k2-sdk/bases"	
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
	"$module/ui"
	{{range $i := $componentList}}
	{{ $i }} "{{ $module }}/components/{{ $i }}"
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
			{{range $i := $componentList}}
			{{ $i }}.NewComponent,
			{{ end }}
		},
	)
}
