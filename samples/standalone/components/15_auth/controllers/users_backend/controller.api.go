package users_backend

import (
	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/system"
)

const cookieKey = "auth"

func (h *Controller) Register(r *gin.RouterGroup) error {
	r.GET("/users/login", h.api_loginGet())
	r.POST("/users/login", h.api_loginPost())
	r.GET("/users/verify", h.api_verify())
	r.GET("/users/logout", h.api_logout())
	return nil
}

// api_verify godoc
// @Summary  Call
// @Schemes
// @Tags files
// @Produce json
// @Success 200  {string} string "OK"
// @Failure 403 {object} string
// @Failure 500 {object} string
// @Router /verify [get]
func (h *Controller) api_verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authCookie, err := ctx.Cookie(cookieKey)
		if err != nil {
			ctx.JSON(403, gin.H{"message": "No auth cookie"})
			return
		}

		configService := h.GetComponent().GetApp().GetKernel().GetService("config").(*config.Service)
		allowedUsers := configService.Get("auth.users").([]any)
		for _, user := range allowedUsers {
			if user.(map[string]any)["username"] == authCookie {
				ctx.Status(200)
				return
			}
		}

		ctx.Status(403)
	}
}

// api_logout godoc
// @Summary  Call
// @Schemes
// @Tags auth
// @Produce json
// @Success 200  {string} string "OK"
// @Failure 403 {object} string
// @Failure 500 {object} string
// @Router /logout [get]
func (h *Controller) api_logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie(cookieKey, "", -1, "/", "", false, true)
		configService := h.GetComponent().GetApp().GetKernel().GetService("config").(*config.Service)
		rootUrl := configService.Get("host.ingress.root_url").(string)

		ctx.Redirect(302, rootUrl)
	}
}

//go:embed pages/login.html
var loginPage []byte

// api_loginGet godoc
// @Summary  Call
// @Schemes
// @Tags auth
// @Produce json
// @Success 200  {string} string "OK"
// @Failure 403 {object} string
// @Failure 500 {object} string
// @Router /login [get]
func (h *Controller) api_loginGet() gin.HandlerFunc {
	configService := h.GetComponent().GetApp().GetKernel().GetService("config").(*config.Service)
	rootUrl := configService.Get("host.ingress.root_url").(string)

	return func(ctx *gin.Context) {
		redirectParam, err := configService.GetAsString("host.ingress.auth.authenticated.redirectParam")
		if err != nil {
			ctx.String(500, "Error getting redirect param: %s", err.Error())
			return
		}

		redirect := ctx.Query(redirectParam)
		if redirect == "" {
			redirect = rootUrl
		}

		page, err := h.buildLoginPage(nil, rootUrl, redirect)
		if err != nil {

			ctx.String(500, "Error building login page %s", err.Error())
			return
		}

		ctx.Data(200, "text/html", []byte(page))
	}
}

// api_loginPost godoc
// @Summary  Call
// @Schemes
// @Tags auth
// @Produce json
// @Success 200  {string} string "OK"
// @Failure 403 {object} string
// @Failure 500 {object} string
// @Router /login [post]
func (h *Controller) api_loginPost() gin.HandlerFunc {
	configService := h.GetComponent().GetApp().GetKernel().GetService("config").(*config.Service)
	rootUrl := configService.Get("host.ingress.root_url").(string)
	allowedUsers := configService.Get("auth.users").([]any)
	return func(ctx *gin.Context) {

		redirect := ctx.PostForm("redirect")
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		for _, user := range allowedUsers {
			if user.(map[string]any)["username"] == username && user.(map[string]any)["password"] == password {
				username := user.(map[string]any)["username"].(string)

				ctx.SetCookie(cookieKey, username, 3600, "/", "", false, true)

				if redirect != "" {
					ctx.Redirect(302, redirect)
				} else {
					ctx.Redirect(302, rootUrl)
				}
				return
			}
		}
		errorMessage := "Invalid credentials"
		page, err := h.buildLoginPage(&errorMessage, rootUrl, redirect)
		if err != nil {
			ctx.String(500, "Error building login page %s", err.Error())
			return
		}

		ctx.Data(403, "text/html", []byte(page))
	}
}

func (h *Controller) buildLoginPage(errorMessage *string, rootUrl string, redirect string) (string, error) {

	return system.UnTemplateWithGoTemplate(string(loginPage), map[string]any{
		"errorMessage": errorMessage,
		"rootUrl":      rootUrl,
		"redirect":     redirect,
	})

}
