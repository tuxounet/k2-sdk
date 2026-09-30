package routes

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

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
	isUnsecure := service.GetKernel().IsUnsecure()

	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if isUnsecure {
			log.TraceF("UNSECURE for path %s", requestPath)
			c.Next()
			return
		}
		if requestPath == verify_url {
			c.Next()
			return
		}

		log.TraceF("check %s", requestPath)

		var matchedPolicy runtimeTypes.IAccessPolicy = ""
		matchedPolicy = matchLongestPrefix(requestPath, authMap)
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
	isUnsecure := service.GetKernel().IsUnsecure()
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if isUnsecure {
			log.TraceF("UNSECURE for path %s", requestPath)
			c.Next()
			return
		}
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

// matchLongestPrefix finds the most specific (longest) matching path prefix
// in the authMap and returns its access policy. If no prefix matches, returns
// an empty string. This guarantees deterministic behavior regardless of the
// authMap iteration order (which is randomized in Go).
func matchLongestPrefix(requestPath string, authMap map[string]runtimeTypes.IAccessPolicy) runtimeTypes.IAccessPolicy {
	var matchedPolicy runtimeTypes.IAccessPolicy = ""
	longestMatch := 0
	requestLen := len(requestPath)
	for pathPrefix, accessPolicy := range authMap {
		prefixLen := len(pathPrefix)
		// Standard prefix match: request starts with the prefix
		standardMatch := requestLen >= prefixLen && requestPath[:prefixLen] == pathPrefix
		// Suffix match: request + "/" equals a trailing-slash prefix (e.g. request "/api/hello"
		// matches the prefix "/api/hello/" which was registered by a controller)
		suffixMatch := requestLen+1 == prefixLen && strings.HasSuffix(pathPrefix, "/") &&
			requestPath+"/" == pathPrefix
		if standardMatch || suffixMatch {
			if prefixLen > longestMatch {
				matchedPolicy = accessPolicy
				longestMatch = prefixLen
			}
		}
	}
	return matchedPolicy
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

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			// ErrAbortHandler is normal for SSE/long-polling disconnections — log nothing
			// For other errors, log them
			if err != http.ErrAbortHandler {
				log.Printf("proxy error: %v", err)
			}
		}

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

		// Recover from http.ErrAbortHandler panics caused by client disconnections
		// during SSE/long-polling. Let other panics propagate to gin's recovery.
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					// client disconnected — this is normal, suppress silently
					return
				}
				panic(err)
			}
		}()

		proxy.ServeHTTP(ctx.Writer, ctx.Request)

	}
}
