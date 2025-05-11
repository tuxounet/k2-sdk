package sample

import (
	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/samples/standalone/components/20_sample_app/controllers/sample2"
	"github.com/tuxounet/k2-sdk/types"
)

type SampleController struct {
	bases.BaseAppController
}

func NewSampleController(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "sample", 1, nil, types.AccessPolicyAuthenticated)
	return &SampleController{
		base,
	}
}

func (h *SampleController) Register(r *gin.RouterGroup) error {

	r.GET("/sayHello", h.api_hello())
	return nil
}

// api_hello godoc
// @Summary  Hello, world!
// @Schemes
// @Tags hello
// @Produce json
// @Success 200  {string} string "OK"
// @Router /sayHello [get]
func (h *SampleController) api_hello() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		msg := h.GetData("greeting")
		if msg == "" {
			msg = "Hello, world!"
		}

		ctrl := h.GetComponent().GetController("sample2").(*sample2.Sample2Controller)
		ctrl.DoIt()

		ctx.JSON(200, gin.H{
			"message": msg,
		})
	}
}
