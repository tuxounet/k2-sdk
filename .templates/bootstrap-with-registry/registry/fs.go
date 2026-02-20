package registry

import (
	"embed"
)

//go:embed dist/*
var RegistryFS embed.FS
