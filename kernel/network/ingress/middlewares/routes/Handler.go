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
		router.Any(fmt.Sprintf("%s/*proxyPath", registration.IngressPath), ensureAuthLevel(service, registration), performProxyRequest(registration))
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
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = ctx.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = ctx.Param("proxyPath")
			req.Header.Set("X-Forwarded-Host", ctx.Request.Host)
			if ctx.Request.TLS == nil {
				req.Header.Set("X-Forwarded-Proto", "http")
			} else {
				req.Header.Set("X-Forwarded-Proto", "https")
			}
			req.Header.Set("X-Forwarded-For", ctx.Request.RemoteAddr)

		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)

		//r.Any(fmt.Sprintf("%s/*proxyPath", ing.Path), func(ctx *gin.Context) {
		// 	localPort, err := computeContainers.LookupIngressPort(b.GetName(), definition.Order, ing.Path)
		// 	if err != nil {
		// 		b.GetLogger().ErrorF("Failed to lookup http port: %s", err)
		// 		ctx.String(500, "Failed to lookup http port")
		// 		return
		// 	}
		// 	if localPort == nil {
		// 		b.GetLogger().ErrorF("No local port found for %s", ing.Path)
		// 		ctx.String(404, "No local port found")
		// 		return
		// 	}

		// 	targetUrl := fmt.Sprintf("http://%s:%d", hostAddr, localPort.LocalPort)
		// 	remote, err := url.Parse(targetUrl)
		// 	if err != nil {
		// 		b.GetLogger().ErrorF("Failed to parse target url: %s", err)
		// 		ctx.String(500, "Failed to parse target url")
		// 		return
		// 	}
		// 	targetPath := ctx.Param("proxyPath")
		// 	proxy := httputil.NewSingleHostReverseProxy(remote)
		// 	proxy.Director = func(req *http.Request) {
		// 		req.Header = ctx.Request.Header
		// 		req.Host = remote.Host
		// 		req.URL.Scheme = remote.Scheme
		// 		req.URL.Host = remote.Host
		// 		req.URL.Path = targetPath
		// 		req.Header.Set("X-Forwarded-Host", ctx.Request.Host)
		// 		if ctx.Request.TLS == nil {
		// 			req.Header.Set("X-Forwarded-Proto", "http")
		// 		} else {
		// 			req.Header.Set("X-Forwarded-Proto", "https")
		// 		}
		// 		req.Header.Set("X-Forwarded-For", ctx.Request.RemoteAddr)

		// 		b.GetLogger().TraceF("Proxying request from %s with headers %v to %s with headers %s", ctx.Request.URL.String(), ctx.Request.Header, req.URL.String(), req.Header)

		// 	}

		// 	proxy.ServeHTTP(ctx.Writer, ctx.Request)

		// })

		// Perform the proxy request

	}
}
