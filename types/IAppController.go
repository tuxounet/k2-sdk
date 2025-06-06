package types

import (
	"embed"

	"github.com/gin-gonic/gin"
)

type IAppController interface {
	ILoggable
	GetName() string
	GetOrder() int
	GetComponent() IAppComponent
	GetConfig() *embed.FS
	GetAccessPolicy() IAccessPolicy
	Init() error
	Register(r *gin.RouterGroup) error
	Start() error
	Stop() error
}

type AppControllerCtor = func(component IAppComponent) IAppController
