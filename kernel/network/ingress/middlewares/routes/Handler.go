package routes

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

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

		// Customize the response by rewriting subpath references in the response body.
		proxy.ModifyResponse = func(resp *http.Response) error {
			// Here, we assume that only certain content types need rewriting.
			contentType := resp.Header.Get("Content-Type")
			if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/plain") {
				// Read the body
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				// Since we've read the body, close the original one.
				resp.Body.Close()

				// Convert to string and perform substitutions.
				// In this example, we simply replace occurrences of href="/ and src="/
				// with an internal link that prepends /subpath.
				bodyStr := string(bodyBytes)
				bodyStr = strings.ReplaceAll(bodyStr, `href="/`, `href="`+ingress.IngressPath+`"`)
				bodyStr = strings.ReplaceAll(bodyStr, `src="/`, `src="`+ingress.IngressPath+`"`)

				// Write the modified body back and update headers accordingly.
				newBody := []byte(bodyStr)
				resp.Body = io.NopCloser(bytes.NewReader(newBody))
				resp.ContentLength = int64(len(newBody))
				resp.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
			}
			return nil
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
