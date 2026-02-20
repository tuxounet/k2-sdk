package registry

import (
	"embed"
)

//go:embed dist/*.so
var RegistryFS embed.FS
