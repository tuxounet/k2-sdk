package bases

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (b *BaseControllerContainer) Register(r *gin.RouterGroup) error {

	provider := b.getComputeContainersProviders()
	definition := b.GetDefinition()
	err := provider.RegisterDefinition(*definition)

	if err != nil {
		b.GetLogger().ErrorF("Failed to register container defintion inside provider: %s", err)
		return err
	}

	if definition == nil || definition.Ingresses == nil || len(definition.Ingresses) == 0 {
		b.GetLogger().TraceF("No ingress defined for container %s, skipping register", b.GetName())
		return nil
	}

	ingresses := definition.Ingresses

	config := b.getConfigService()
	hostAddr, err := config.GetAsString("host.address")
	if err != nil {
		b.GetLogger().ErrorF("Failed to get host address: %s", err)
		return err
	}

	computeContainers := b.getComputeContainersProviders()

	for _, ing := range ingresses {

		r.Any(fmt.Sprintf("%s/*proxyPath", ing.Path), func(ctx *gin.Context) {
			localPort, err := computeContainers.LookupIngressPort(b.GetName(), definition.Order, ing.Path)
			if err != nil {
				b.GetLogger().ErrorF("Failed to lookup http port: %s", err)
				ctx.String(500, "Failed to lookup http port")
				return
			}
			if localPort == nil {
				b.GetLogger().ErrorF("No local port found for %s", ing.Path)
				ctx.String(404, "No local port found")
				return
			}

			targetUrl := fmt.Sprintf("http://%s:%d", hostAddr, localPort.LocalPort)
			remote, err := url.Parse(targetUrl)
			if err != nil {
				b.GetLogger().ErrorF("Failed to parse target url: %s", err)
				ctx.String(500, "Failed to parse target url")
				return
			}
			targetPath := ctx.Param("proxyPath")
			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.Director = func(req *http.Request) {
				req.Header = ctx.Request.Header
				req.Host = remote.Host
				req.URL.Scheme = remote.Scheme
				req.URL.Host = remote.Host
				req.URL.Path = targetPath
				req.Header.Set("X-Forwarded-Host", ctx.Request.Host)
				if ctx.Request.TLS == nil {
					req.Header.Set("X-Forwarded-Proto", "http")
				} else {
					req.Header.Set("X-Forwarded-Proto", "https")
				}
				req.Header.Set("X-Forwarded-For", ctx.Request.RemoteAddr)

				b.GetLogger().TraceF("Proxying request from %s with headers %v to %s with headers %s", ctx.Request.URL.String(), ctx.Request.Header, req.URL.String(), req.Header)

			}

			proxy.ServeHTTP(ctx.Writer, ctx.Request)

		})
	}

	return nil
}
