package apps

import (
	"strings"

	"github.com/gin-gonic/gin"
	runtimeBases "github.com/tuxounet/k2-sdk/bases"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	runtimeBases.BaseAppController
}

func NewController(component runtimeTypes.IAppComponent) runtimeTypes.IAppController {
	base := runtimeBases.NewBaseAppController(component, "apps", 2, nil, runtimeTypes.AccessPolicyPublic)
	return &Controller{base}
}

func (h *Controller) Register(r *gin.RouterGroup) error {
	r.GET("/list", h.api_listApps())
	return nil
}

type AppEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// api_listApps godoc
// @Summary  List available external applications
// @Schemes
// @Tags apps
// @Produce json
// @Success 200  {array} AppEntry
// @Router /site/apps/list [get]
func (h *Controller) api_listApps() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		app := h.GetComponent().GetApp()
		var result []AppEntry

		externals := app.GetExternals()
		if externals == nil {
			ctx.JSON(200, []AppEntry{})
			return
		}

		entries, err := externals.ReadDir("dist")
		if err != nil {
			ctx.JSON(200, []AppEntry{})
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !strings.HasSuffix(name, ".so") {
				continue
			}
			componentName := strings.TrimSuffix(name, ".so")

			// Verify the component is actually loaded and has a UI
			component := app.GetComponent(componentName)
			if component != nil && component.GetUI() != nil {
				result = append(result, AppEntry{
					Name: componentName,
					URL:  "/" + componentName + "/ui/",
				})
			}
		}

		ctx.JSON(200, result)
	}
}
