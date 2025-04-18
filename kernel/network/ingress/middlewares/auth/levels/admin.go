package levels

import (
	"net/http"

	"github.com/tuxounet/k2-sdk/types"
)

func AllowAccessLevelAdmin(_ *http.Request, _ types.ILogger) bool {
	//not implemented

	return false
}
