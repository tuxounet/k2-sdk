package types

import (
	"embed"

	"github.com/swaggo/swag"
)

type IApp interface {
	ILoggable
	GetName() string
	GetVersion() string
	GetDocs() *swag.Spec
	GetUI() *embed.FS
	GetExternals() *embed.FS
	GetConfig() *embed.FS
	SetLogger(logger ILogger)
	GetComponents() []IAppComponent
	GetComponent(name string) IAppComponent
	AddComponent(ctor AppComponentCtor)
	GetKernel() IKernel
	SetKernel(kernel IKernel)
}
