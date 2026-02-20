package temps

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) Register(r *gin.RouterGroup) error {
	r.GET("/time", c.api_getTime)
	return nil
}

// api_getTime godoc
//
//	@Summary		Retourne l'heure actuelle du serveur
//	@Description	Retourne l'heure du serveur au format français avec toutes les informations nécessaires pour l'affichage analogique et numérique
//	@Schemes
//	@Tags		Horloge
//	@Produce	json
//	@Success	200	{object}	TimeResponse	"Heure actuelle"
//	@Router		/horloge/temps/time [get]
func (c *Controller) api_getTime(ctx *gin.Context) {
	response := c.GetCurrentTime()
	ctx.JSON(200, response)
}
