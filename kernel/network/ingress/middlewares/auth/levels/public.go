package levels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/types"
)

func isPublicEnabled(configService *config.Service) bool {
	public_enabled, err := configService.GetAsBool("host.ingress.auth.public.enabled")
	if err != nil {
		return false
	}
	return public_enabled
}

func AllowAccessLevelPublic(req *http.Request, log types.ILogger, configService *config.Service) bool {

	requestUrl := req.URL.Path

	if !isPublicEnabled(configService) {
		log.InfoF("Public access is disabled for path: %s", requestUrl)
		return false
	}
	return true
}

func BlockAccessLevelPublic(req *http.Request, log types.ILogger, configService *config.Service, c *gin.Context) {
	requestUrl := req.URL.Path
	if !isPublicEnabled(configService) {
		log.InfoF("Public access is disabled for path: %s", requestUrl)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	log.WarnF("Access denied to public level: %s", requestUrl)
	c.AbortWithStatus(http.StatusForbidden)
}
