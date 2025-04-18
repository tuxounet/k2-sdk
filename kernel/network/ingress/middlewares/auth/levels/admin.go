package levels

import (
	"net/http"

	"github.com/tuxounet/k2-sdk/types"
)

func AllowAccessLevelAdmin(req *http.Request, log types.ILogger) bool {
	requestUrl := req.URL.Path

	log.ErrorF("Access denied to admin level: %s", requestUrl)

	return false
}
