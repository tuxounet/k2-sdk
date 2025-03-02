package bases

import (
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/types"
)

type BaseAppController struct {
	name      string
	order     int
	component types.IAppComponent
	log       types.ILogger
	config    *embed.FS
	data      map[string]interface{}
}

func NewBaseAppController(component types.IAppComponent, name string, order int, config *embed.FS) BaseAppController {

	log := component.GetLogger().CreateSubLogger(name)
	base := BaseAppController{
		name:      name,
		component: component,
		order:     order,
		config:    config,
		log:       log,
		data:      make(map[string]interface{}),
	}

	return base
}

func (b *BaseAppController) GetName() string {
	return b.name
}

func (b *BaseAppController) GetOrder() int {
	return b.order
}

func (b *BaseAppController) GetComponent() types.IAppComponent {
	return b.component
}

func (b *BaseAppController) GetLogger() types.ILogger {
	return b.log
}

func (b *BaseAppController) Init() error {
	return nil
}

func (b *BaseAppController) Register(router *gin.RouterGroup) error {
	return nil
}

func (b *BaseAppController) Start() error {
	return nil
}
func (b *BaseAppController) Stop() error {
	return nil
}

func (b *BaseAppController) GetConfig() *embed.FS {
	return b.config
}
func (a *BaseAppController) GetData(key string) interface{} {
	return a.data[key]
}

func (a *BaseAppController) SetData(key string, value interface{}) {

	a.data[key] = value
}
