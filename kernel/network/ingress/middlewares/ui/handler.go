package ui

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/ui/files"
	"github.com/tuxounet/k2-sdk/types"
)

func RegisterApp(app types.IApp, rootUri *url.URL, router *gin.RouterGroup) error {

	log := app.GetLogger().CreateSubLogger("ui")
	baseUri := router.BasePath()
	fs := app.GetUI()
	if fs == nil {
		log.Warn("No UI folder found for app")
		return nil
	}
	fsPrefix := "dist"
	uiPathSuffix := "ui/"

	mw, err := files.ServeAppMiddleware(app, log, baseUri, fs, fsPrefix)
	if err != nil {
		log.ErrorF("Failed to create serve app middleware: %s", err.Error())
		return err
	}
	router.GET(uiPathSuffix+"/*any", mw)

	//Initial redirection
	router.GET("/", func(c *gin.Context) {
		target := fmt.Sprintf("%s%s/", baseUri, uiPathSuffix)
		c.Redirect(302, target)
	})

	return nil
}

func RegisterComponent(component types.IAppComponent, rootUri *url.URL, baseUri string, router *gin.RouterGroup) error {

	log := component.GetLogger().CreateSubLogger("ui")

	fs := component.GetUI()
	if fs == nil {
		log.Warn("No UI folder found for component")
		return nil
	}
	fsPrefix := "dist"
	uiPathSuffix := "ui/"

	mw, err := files.ServeComponentMiddleware(component, log, baseUri, fs, fsPrefix)
	if err != nil {
		log.ErrorF("Failed to create serve component middleware: %s", err.Error())
		return err
	}
	router.GET(uiPathSuffix+"/*any", mw)

	//Initial redirection
	router.GET("/", func(c *gin.Context) {
		target := fmt.Sprintf("%s/%s/", baseUri, uiPathSuffix)
		c.Redirect(302, target)
	})

	return nil
}
