# k2-sdk

## Install

```sh
go get -v  github.com/tuxounet/k2-sdk
```

## Bootstrap a new App

### Makefile

```Makefile
APP_NAME := Sample
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),v${GIT_BRANCH})
GO_PATH:=$(shell go env GOPATH)
VERSION_FILE := ./app/version.go
init:
	go mod tidy

write-version:
	echo "package app" > ${VERSION_FILE}
	echo "" >> ${VERSION_FILE}
	echo "const (" >> ${VERSION_FILE}
	echo "    AppName    = \"${APP_NAME}\"" >> ${VERSION_FILE}
	echo "    AppVersion = \"${VERSION}\"" >> ${VERSION_FILE}
	echo ")" >> ${VERSION_FILE}

prepare: init write-version
	go install github.com/swaggo/swag/cmd/swag@v1.16.4
	${GO_PATH}/bin/swag init --instanceName ${APP_NAME}
	${GO_PATH}/bin/swag fmt

run: prepare
	go run ./main.go
```

### Entrypoint (main.go)

```golang
package main

import (
	"myApp/app"
	runtime "github.com/tuxounet/k2-sdk"
)

// @title			Sample
// @version		0.0
// @description	This is the API for Sample App
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	GPL-3.0
// @license.url	http://www.gnu.org/licenses/gpl-3.0.html
func main() {
	runtime.HostSingleApp(app.NewApp())
}
```

### Sample controller (controllers/hello/controller.go)

```golang
package hello

import (
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
	bases.BaseController
}

func NewHelloController(app types.IApp) types.IController {
	base := bases.NewBaseController(app, "hello")
	return &HelloController{
		base,
	}
}

func (h *HelloController) Register(r *gin.RouterGroup) error {

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
func (h *HelloController) api_hello() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.GetLogger().Info("Something said hello!")
		ctx.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	}
}

```

### App Definition

```golang

package app

import (
	"myApp/controllers/hello"
	"myApp/docs"
	"myApp/ui"
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

func NewApp() types.IApp {
	return bases.NewBaseApp(
		AppName,
		AppVersion,
		docs.SwaggerInfoHello,
		&ui.Dist,
		[]types.ControllerCtor{
			hello.NewHelloController,
		},
	)
}

```
