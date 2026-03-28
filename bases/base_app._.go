package bases

import (
	"embed"
	"strings"

	"github.com/swaggo/swag"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type GenericApp struct {
	BaseApp
}

type BaseApp struct {
	name               string
	version            string
	docs               *swag.Spec
	ui                 *embed.FS
	config             *embed.FS
	log                runtimeTypes.ILogger
	kernel             runtimeTypes.IKernel
	componentsCtors    []runtimeTypes.AppComponentCtor
	components         []runtimeTypes.IAppComponent
	externalComponents *embed.FS
}

func NewBaseApp(name string, version string, docs *swag.Spec, ui *embed.FS, config *embed.FS, componentsCtors []runtimeTypes.AppComponentCtor, externalComponents *embed.FS) runtimeTypes.IApp {

	name = strings.TrimSpace(name)
	version = strings.TrimSpace(version)

	if docs != nil {
		docs = &swag.Spec{}

		if name != "" {
			docs.Title = name
		}
		if version != "" {
			docs.Version = version
		}
		docs.BasePath = "/"
	}

	base := BaseApp{
		name:               name,
		version:            version,
		docs:               docs,
		ui:                 ui,
		config:             config,
		componentsCtors:    componentsCtors,
		externalComponents: externalComponents,
	}

	return &GenericApp{base}

}
