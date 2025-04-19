package auth

import (
	"net/http"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/auth/levels"
	"github.com/tuxounet/k2-sdk/types"
)

func Register(service types.IKernelService, router *gin.RouterGroup) error {

	log := service.GetLogger().CreateSubLogger("auth")

	configService := service.GetKernel().GetService(config.ServiceKey).(*config.Service)

	default_access, err := configService.GetAsString("host.ingress.auth.defaultAccess")
	if err != nil {
		log.ErrorF("Failed to get defaultAccess: %s", err.Error())
		return err
	}
	defaut_policy := types.IAccessPolicy(default_access)

	router.Use(func(c *gin.Context) {

		accessLevel := accessLevelEvaluation(service, c.Request, defaut_policy, log)
		switch accessLevel {
		case types.AccessPolicyPublic:
			if levels.AllowAccessLevelPublic(c.Request, log, configService) {
				c.Next()
				return
			} else {
				levels.BlockAccessLevelPublic(c.Request, log, configService, c)
				return
			}

		case types.AccessPolicyAuthenticated:
			if levels.AllowAccessLevelAuthenticated(c.Request, log, configService) {
				c.Next()
				return
			} else {
				levels.RedirectAccessLevelAuthenticatedLogin(c.Request, log, configService, c)
				return
			}
		case types.AccessPolicyAdmin:
			if levels.AllowAccessLevelAdmin(c.Request, log, configService) {
				c.Next()
				return
			} else {
				levels.RedirectAccessLevelAdminLogin(c.Request, log, configService, c)
				return
			}

		default:
			log.ErrorF("Unknown access level: %s", accessLevel)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

	})

	return nil
}

func accessLevelEvaluation(service types.IKernelService, request *http.Request, default_access types.IAccessPolicy, log types.ILogger) types.IAccessPolicy {
	requestPath := request.URL.Path

	cleanSegments := make([]string, 0)
	segmets := strings.Split(requestPath, "/")
	for _, seg := range segmets {
		if seg == "" {
			continue
		}
		cleanSegments = append(cleanSegments, seg)
	}
	if len(cleanSegments) == 0 {
		log.DebugF("Access level for %s is %s", requestPath, default_access)
		return default_access
	}
	first := cleanSegments[0]
	allowedSegments := []string{"ui", "cdn"}
	if slices.Contains(allowedSegments, first) {
		log.DebugF("Access level for %s is %s", requestPath, default_access)
		return default_access
	}

	app := service.GetKernel().GetApp()
	components := app.GetComponents()
	for _, component := range components {
		if component.GetName() == first {
			policy := component.GetAccessPolicy()
			log.DebugF("Access level for %s is %s", requestPath, policy)
			return policy
		}
	}

	return default_access
}
