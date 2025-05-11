package sample2

import (
	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/kernel/profile"
	"github.com/tuxounet/k2-sdk/types"
)

type Sample2Controller struct {
	bases.BaseAppController
}

func NewSample2Controller(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "sample2", 3, nil, types.AccessPolicyAuthenticated)
	return &Sample2Controller{
		base,
	}
}

func (h *Sample2Controller) DoIt() error {

	return h.GetLogger().Scope("do it", func(log types.ILogger) error {
		profile := h.GetComponent().GetApp().GetKernel().GetService(profile.ServiceKey).(*profile.ProfileService)
		currentProfile, err := profile.GetPublicProfile()
		if err != nil {
			return err
		}

		dataDir, err := profile.GetUserDirectory()
		if err != nil {
			return err
		}

		log.InfoF("Current Profile: %v - %s", currentProfile, dataDir)
		return nil
	})

}

func (h *Sample2Controller) Register(r *gin.RouterGroup) error {

	r.GET("/call", h.api_call())
	return nil
}

// api_call godoc
// @Summary  Call
// @Schemes
// @Tags hello
// @Produce json
// @Success 200  {string} string "OK"
// @Router /call [get]
func (h *Sample2Controller) api_call() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h.DoIt()
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "done",
		})
	}
}
