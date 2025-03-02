package health

import (
	"github.com/gin-gonic/gin"
	runtimeBases "{{ .sdk_module }}/bases"
	runtimeTypes "{{ .sdk_module }}/types"
)

type HealthController struct {
	runtimeBases.BaseAppController
}

func NewController(component runtimeTypes.IAppComponent) runtimeTypes.IAppController {
	base := runtimeBases.NewBaseAppController(component, "health", 1, nil)
	return &HealthController{
		base,
	}
}

func (h *HealthController) Register(r *gin.RouterGroup) error {

	r.GET("/health", h.api_health())
	return nil
}

// api_health godoc
// @Summary  Hello, world!
// @Schemes
// @Tags hello
// @Produce json
// @Success 200  {string} string "OK"
// @Router /sayHello [get]
func (h *HealthController) api_health() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, "OK")
	}
}
