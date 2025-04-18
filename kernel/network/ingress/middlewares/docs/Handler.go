package docs

import (
	_ "embed"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"github.com/tuxounet/k2-sdk/types"
)

func Register(service types.IKernelService, docs *swag.Spec, router *gin.RouterGroup) error {

	log := service.GetLogger().CreateSubLogger("docs")

	baseRoute := router.BasePath()

	if !strings.HasSuffix(baseRoute, "/") {
		baseRoute += "/"
	}

	docsPathSuffix := "docs/"
	router.GET(docsPathSuffix+"*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL(baseRoute+"openapi.json")),
	)

	router.GET(baseRoute+"openapi.json", func(c *gin.Context) {
		swagger := swag.GetSwagger(docs.InfoInstanceName)
		if swagger == nil {
			log.ErrorF("Error parsing template for %s", docs.InfoInstanceName)
			c.String(500, "Error parsing template")
			return
		}

		body := swagger.ReadDoc()
		c.Data(200, "application/json", []byte(body))
	})

	return nil

}
