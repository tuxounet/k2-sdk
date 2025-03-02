package bases

import (
	"embed"

	"github.com/swaggo/swag"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type GenericApp struct {
	BaseApp
}

type BaseApp struct {
	name            string
	version         string
	docs            *swag.Spec
	ui              *embed.FS
	config          *embed.FS
	log             runtimeTypes.ILogger
	kernel          runtimeTypes.IKernel
	componentsCtors []runtimeTypes.AppComponentCtor
	components      []runtimeTypes.IAppComponent
}

func NewBaseApp(name string, version string, docs *swag.Spec, ui *embed.FS, config *embed.FS, componentsCtors []runtimeTypes.AppComponentCtor) runtimeTypes.IApp {

	if docs == nil {
		docs = &swag.Spec{}
	}

	if name != "" {
		docs.Title = name
	}
	if version != "" {
		docs.Version = version
	}

	docs.BasePath = "/"

	base := BaseApp{
		name:            docs.Title,
		version:         docs.Version,
		docs:            docs,
		ui:              ui,
		config:          config,
		componentsCtors: componentsCtors,
	}

	return &GenericApp{base}

}
