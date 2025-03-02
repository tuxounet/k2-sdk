package bases

import (
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type BaseAppComponent struct {
	name             string
	order            int
	app              runtimeTypes.IApp
	docs             *swag.Spec
	ui               *embed.FS
	config           *embed.FS
	log              runtimeTypes.ILogger
	controllersCtors []runtimeTypes.AppControllerCtor
	controllers      []runtimeTypes.IAppController
	accessPolicy     runtimeTypes.IAccessPolicy
}

func NewBaseAppComponent(
	app runtimeTypes.IApp,
	name string,
	order int,
	docs *swag.Spec,
	ui *embed.FS,
	config *embed.FS,
	accessPolicy runtimeTypes.IAccessPolicy,
	controllersCtors []runtimeTypes.AppControllerCtor) runtimeTypes.IAppComponent {

	log := app.GetLogger().CreateSubLogger(name)
	base := &BaseAppComponent{
		name:             name,
		app:              app,
		order:            order,
		config:           config,
		log:              log,
		docs:             docs,
		ui:               ui,
		accessPolicy:     accessPolicy,
		controllersCtors: controllersCtors,
	}

	return base
}

func (b *BaseAppComponent) GetName() string {
	return b.name
}

func (b *BaseAppComponent) GetOrder() int {
	return b.order
}

func (b *BaseAppComponent) GetApp() runtimeTypes.IApp {
	return b.app
}

func (b *BaseAppComponent) GetLogger() runtimeTypes.ILogger {
	return b.log
}
func (a *BaseAppComponent) GetDocs() *swag.Spec {

	return a.docs
}

func (a *BaseAppComponent) GetUI() *embed.FS {
	return a.ui
}

func (a *BaseAppComponent) GetAccessPolicy() runtimeTypes.IAccessPolicy {
	return a.accessPolicy
}

func (a *BaseAppComponent) GetControllers() []runtimeTypes.IAppController {
	if a.controllers == nil {

		a.controllers = make([]runtimeTypes.IAppController, 0)
		for _, ctor := range a.controllersCtors {
			component := ctor(a)
			a.controllers = append(a.controllers, component)
		}
	}

	return a.controllers
}

func (a *BaseAppComponent) GetController(name string) runtimeTypes.IAppController {
	for _, controller := range a.GetControllers() {
		if controller.GetName() == name {
			return controller
		}
	}

	return nil
}

func (b *BaseAppComponent) Init() error {
	return nil
}

func (b *BaseAppComponent) Register(router *gin.RouterGroup) error {
	return nil
}

func (b *BaseAppComponent) Start() error {
	return nil
}
func (b *BaseAppComponent) Stop() error {
	return nil
}

func (b *BaseAppComponent) GetConfig() *embed.FS {
	return b.config
}
