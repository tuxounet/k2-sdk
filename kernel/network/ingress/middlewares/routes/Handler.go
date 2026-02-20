package routes

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/routes/access"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

func RegisterIngresses(service runtimeTypes.IKernelService, router *gin.RouterGroup, ingresses []types.IngressDefinition) error {

	for _, registration := range ingresses {
		handler := performProxyRequest(registration)
		if registration.CustomHandler != nil {
			handler = registration.CustomHandler(registration)
		}
		router.Any(fmt.Sprintf("%s/*proxyPath", registration.IngressPath), ensureAuthLevel(service, registration), handler)
	}

	return nil
}

func EnsureAuthLevelMiddleware(service runtimeTypes.IKernelService, parentLog runtimeTypes.ILogger, authMap map[string]runtimeTypes.IAccessPolicy, defaultAccess runtimeTypes.IAccessPolicy) gin.HandlerFunc {
	log := parentLog.CreateSubLogger("auth")
	configService := service.GetKernel().GetService(config.ServiceKey).(*config.Service)

	login_url, _ := configService.GetAsStringOrDefault("host.ingress.auth.authenticated.login_url", "")
	verify_url, _ := configService.GetAsStringOrDefault("host.ingress.auth.authenticated.verify_url", "")

	authMap[login_url] = runtimeTypes.AccessPolicyPublic
	authMap[verify_url] = runtimeTypes.AccessPolicyPublic

	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if requestPath == verify_url {
			c.Next()
			return
		}
		log.TraceF("check %s", requestPath)

		var matchedPolicy runtimeTypes.IAccessPolicy = ""
		for pathPrefix, accessPolicy := range authMap {
			if len(requestPath) >= len(pathPrefix) && requestPath[0:len(pathPrefix)] == pathPrefix {
				matchedPolicy = accessPolicy
			}
		}
		if matchedPolicy == "" {
			matchedPolicy = defaultAccess
		}
		log.TraceF("matched policy %s for path %s", matchedPolicy, requestPath)

		switch matchedPolicy {
		case runtimeTypes.AccessPolicyPublic:
			if access.AllowAccessLevelPublic(c.Request, log, configService) {
				c.Next()
				return
			} else {
				access.BlockAccessLevelPublic(c.Request, log, configService, c)
				return
			}

		case runtimeTypes.AccessPolicyAuthenticated:
			if access.AllowAccessLevelAuthenticated(c.Request, log, configService) {
				c.Next()
				return
			} else {
				access.RedirectAccessLevelAuthenticatedLogin(c.Request, log, configService, c)
				return
			}

		default:
			log.ErrorF("Unknown access level: %s", matchedPolicy)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
}

func ensureAuthLevel(service runtimeTypes.IKernelService, ingress types.IngressDefinition) gin.HandlerFunc {
	configService := service.GetKernel().GetService(config.ServiceKey).(*config.Service)
	log := service.GetLogger().CreateSubLogger("auth")
	return func(c *gin.Context) {

		switch ingress.AccessPolicy {
		case runtimeTypes.AccessPolicyPublic:
			if access.AllowAccessLevelPublic(c.Request, log, configService) {
				c.Next()
				return
			} else {
				access.BlockAccessLevelPublic(c.Request, log, configService, c)
				return
			}

		case runtimeTypes.AccessPolicyAuthenticated:
			if access.AllowAccessLevelAuthenticated(c.Request, log, configService) {
				c.Next()
				return
			} else {
				access.RedirectAccessLevelAuthenticatedLogin(c.Request, log, configService, c)
				return
			}

		default:
			log.ErrorF("Unknown access level: %s", ingress.AccessPolicy)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

	}

}

func performProxyRequest(ingress types.IngressDefinition) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		targetUri := fmt.Sprintf("http://%s:%d", ingress.ServiceHost, ingress.ServicePort)

		remote, err := url.Parse(targetUri)
		if err != nil {

			ctx.String(500, "Failed to parse target url")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		queryPath := ctx.Request.URL.Path
		if ingress.RewritePath != nil {
			queryPath = *ingress.RewritePath + ctx.Param("proxyPath")
		}

		proxy.Director = func(req *http.Request) {
			req.Header = ctx.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = queryPath
			req.Host = remote.Host
			req.Header.Set("X-Forwarded-Host", ctx.Request.Host)
			if ctx.Request.TLS == nil {
				req.Header.Set("X-Forwarded-Proto", "http")
			} else {
				req.Header.Set("X-Forwarded-Proto", "https")
			}
			req.Header.Set("X-Forwarded-For", ctx.Request.RemoteAddr)

		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)

	}
}
