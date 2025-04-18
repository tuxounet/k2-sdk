package auth

import (
	"fmt"
	"net/http"
	"net/url"
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

	login_url, err := configService.GetAsString("host.ingress.auth.loginUrl")
	if err != nil {
		log.ErrorF("Failed to get loginUrl: %s", err.Error())
		return err
	}

	loginUrl, err := url.Parse(login_url)
	if err != nil {
		log.ErrorF("Failed to parse loginUrl: %s", err.Error())
		return err
	}

	verify_url, err := configService.GetAsString("host.ingress.auth.verifyUrl")
	if err != nil {
		log.ErrorF("Failed to get verifyUrl: %s", err.Error())
		return err
	}
	redirect_param, err := configService.GetAsString("host.ingress.auth.redirectParam")
	if err != nil {
		log.ErrorF("Failed to get redirectParam: %s", err.Error())
		return err
	}

	if !strings.HasPrefix(verify_url, "http") {
		rootUrl, err := configService.GetAsString("host.ingress.rootUrl")
		if err != nil {
			log.ErrorF("Failed to get rootUrl: %s", err.Error())
			return err
		}
		verify_url = fmt.Sprintf("%s%s", rootUrl, verify_url)

	}

	router.Use(func(c *gin.Context) {

		accessLevel := accessLevelEvaluation(service, c.Request, defaut_policy, log)
		switch accessLevel {
		case types.AccessPolicyPublic:
			if levels.AllowAccessLevelPublic(c.Request, log) {
				c.Next()
				return
			} else {
				levels.BlockAccessLevelPublic(c.Request, log, c)
				return
			}

		case types.AccessPolicyAuthenticated:
			if levels.AllowAccessLevelAuthenticated(c.Request, log, verify_url, redirect_param) {
				c.Next()
				return
			} else {
				q := loginUrl.Query()
				q.Add(redirect_param, c.Request.URL.Path)
				loginUrl.RawQuery = q.Encode()

				c.Redirect(http.StatusTemporaryRedirect, loginUrl.String())
				c.Abort()
				return
			}
		case types.AccessPolicyAdmin:
			if levels.AllowAccessLevelAdmin(c.Request, log) {
				c.Next()
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
