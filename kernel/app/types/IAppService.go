package types

import "github.com/tuxounet/k2-sdk/types"

type UnTemplateBody struct {
	Name       string
	Version    string
	RootUrl    string
	BasePath   string
	UIBasePath string
}
type IAppsService interface {
	UnTemplate(app types.IApp, templated string) (string, error)
}
