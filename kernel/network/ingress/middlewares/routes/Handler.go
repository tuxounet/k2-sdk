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

func Register(service runtimeTypes.IKernelService, router *gin.RouterGroup, registrations []types.IngressDefinition) error {

	for _, registration := range registrations {
		handler := registration.CustomHandler
		if handler == nil {
			handler = performProxyRequest(registration)
		}
		router.Any(fmt.Sprintf("%s/*proxyPath", registration.IngressPath), ensureAuthLevel(service, registration), handler)
	}

	return nil
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
		case runtimeTypes.AccessPolicyAdmin:
			if access.AllowAccessLevelAdmin(c.Request, log, configService) {
				c.Next()
				return
			} else {
				access.RedirectAccessLevelAdminLogin(c.Request, log, configService, c)
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

		requestPath := ctx.Param("proxyPath")

		proxy := httputil.NewSingleHostReverseProxy(remote)

		proxy.Director = func(req *http.Request) {
			req.Header = ctx.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = requestPath
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
