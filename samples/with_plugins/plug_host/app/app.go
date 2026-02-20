package app

import (
	"embed"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/samples/with_plugins/plug_host/registry"
	"github.com/tuxounet/k2-sdk/types"
)

//go:embed config/*
var conf embed.FS

func NewApp() types.IApp {
	return bases.NewBaseApp(
		AppName,
		AppVersion,
		nil,
		nil,
		&conf,
		[]types.AppComponentCtor{},
		&registry.RegistryFS,
	)
}
