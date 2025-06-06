package files

import (
	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores"
	"github.com/tuxounet/k2-sdk/types"
)

type FilesController struct {
	bases.BaseAppController
}

func NewFilesController(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "files", 0, nil, types.AccessPolicyAuthenticated)
	return &FilesController{
		base,
	}
}

func (h *FilesController) DoIt() error {

	stores := h.GetComponent().GetApp().GetKernel().GetService(stores.ServiceKey).(*stores.Service)
	store, err := stores.GetStore("local")
	if err != nil {
		return err
	}

	found, err := store.Exists("test")
	if err != nil {
		return err
	}

	if !found {
		err = store.WriteObject("test", []byte("hello world"))
		if err != nil {
			return err
		}
	}

	body, err := store.ReadObject("test")
	if err != nil {
		return err
	}

	h.GetLogger().InfoF("Read: %s", string(body))

	return nil
}

func (h *FilesController) Register(r *gin.RouterGroup) error {

	r.GET("/files/call", h.api_call())
	return nil
}

// api_call godoc
// @Summary  Call
// @Schemes
// @Tags files
// @Produce json
// @Success 200  {string} string "OK"
// @Router /files/call [get]
func (h *FilesController) api_call() gin.HandlerFunc {
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
