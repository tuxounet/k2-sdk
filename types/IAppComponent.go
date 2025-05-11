package types

import (
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

type IAppComponent interface {
	ILoggable
	GetName() string
	GetOrder() int
	GetApp() IApp
	GetConfig() *embed.FS
	GetDocs() *swag.Spec
	GetUI() *embed.FS

	GetControllers() []IAppController
	GetController(name string) IAppController

	Init() error
	Register(r *gin.RouterGroup) error
	Start() error
	Stop() error
}

type AppComponentCtor = func(app IApp) IAppComponent
