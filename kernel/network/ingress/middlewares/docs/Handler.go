package docs

import (
	_ "embed"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"github.com/tuxounet/k2-sdk/kernel/config"
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

	configService := service.GetKernel().GetService(config.ServiceKey).(*config.Service)
	rootUrl, err := configService.GetAsString("host.ingress.rootUrl")
	if err != nil {
		log.ErrorF("Error getting rootUrl: %v", err)
		return err
	}
	paredUrl, err := url.Parse(rootUrl)
	if err != nil {
		log.ErrorF("Error parsing rootUrl: %v", err)
		return err
	}

	router.GET("/openapi.json", func(c *gin.Context) {
		docs.Schemes = []string{paredUrl.Scheme}
		docs.Host = paredUrl.Host
		docs.BasePath = baseRoute
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
