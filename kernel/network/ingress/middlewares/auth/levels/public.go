package levels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/types"
)

func AllowAccessLevelPublic(_ *http.Request, _ types.ILogger) bool {
	return true
}

func BlockAccessLevelPublic(req *http.Request, log types.ILogger, c *gin.Context) {
	requestUrl := req.URL.Path
	log.WarnF("Access denied to public level: %s", requestUrl)
	c.AbortWithStatus(http.StatusForbidden)
}
