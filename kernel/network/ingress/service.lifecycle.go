package ingress

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/auth"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/cdn"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/docs"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/logger"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress/middlewares/ui"
)

func (s *Service) Init() error {

	return nil
}

func (s *Service) Register() error {
	configService := s.getConfigService()

	hostAddr, err := configService.GetAsString("host.ingress.address")
	if err != nil {
		return err
	}

	server := gin.New()
	trustedProxies := []string{hostAddr}
	server.SetTrustedProxies(trustedProxies)

	server.Use(cors.New(cors.Config{
		AllowMethods:  []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Origin, Content-Length, Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))
	s.setServer(server)

	err = logger.Register(s, server)
	if err != nil {
		return err
	}

	rootUrl, err := configService.GetAsString("host.ingress.rootUrl")
	if err != nil {
		s.GetLogger().ErrorF("Failed to get rootUrl: %s", err.Error())
		return err
	}

	rootUri, err := url.Parse(rootUrl)
	if err != nil {
		s.GetLogger().ErrorF("Failed to parse rootUrl: %s", err.Error())
		return err
	}

	kernel := s.GetKernel()
	app := kernel.GetApp()

	baseUrl := rootUri.Path
	if baseUrl == "" {
		baseUrl = "/"
	}

	router := server.Group(baseUrl)

	err = auth.Register(s, router)
	if err != nil {
		s.GetLogger().ErrorF("Failed to register auth: %s", err.Error())
		return err
	}

	components := app.GetComponents()

	appUi := app.GetUI()

	if appUi != nil {
		err := cdn.Register(s, router)
		if err != nil {
			s.GetLogger().ErrorF("Failed to register cdn for app %s: %s", app.GetName(), err.Error())
			return err
		}

		err = ui.RegisterApp(app, rootUri, router.BasePath(), router)
		if err != nil {
			s.GetLogger().ErrorF("Failed to register ui for app %s: %s", app.GetName(), err.Error())
			return err
		}

	}

	for _, component := range components {
		componentRouter := router.Group(component.GetName())

		controllers := component.GetControllers()
		for _, ctrl := range controllers {
			err := ctrl.Register(componentRouter)
			if err != nil {
				s.GetLogger().ErrorF("controller %s in component %s register failed: %s", ctrl.GetName(), component.GetName(), err.Error())
				return err
			}
		}

		componentDocs := component.GetDocs()
		if componentDocs != nil {
			err := docs.Register(s, component, componentRouter)
			if err != nil {
				s.GetLogger().ErrorF("Failed to register docs for component %s: %s", component.GetName(), err.Error())
				return err
			}
		}

		componentUI := component.GetUI()
		if componentUI != nil {
			err := cdn.Register(s, componentRouter)
			if err != nil {
				s.GetLogger().ErrorF("Failed to register cdn for app %s: %s", app.GetName(), err.Error())
				return err
			}

			err = ui.RegisterComponent(component, rootUri, componentRouter.BasePath(), componentRouter)
			if err != nil {
				s.GetLogger().ErrorF("Failed to register ui for app %s: %s", app.GetName(), err.Error())
				return err
			}

		}

	}

	return nil
}

func (s *Service) Listen() error {
	log := s.GetLogger()
	server := s.GetServer()
	if server == nil {
		return errors.New("no server found")
	}

	configService := s.getConfigService()

	hostAddr, err := configService.GetAsString("host.ingress.address")
	if err != nil {
		return err
	}
	hostPort, err := configService.GetAsInt("host.ingress.port")
	if err != nil {
		return err
	}

	rootUrl, err := configService.GetAsString("host.ingress.rootUrl")
	if err != nil {
		return err
	}

	enableTls := configService.Has("host.ingress.tls.port")
	if enableTls {
		//SSL Mode

		listenSSLPort, err := configService.GetAsInt("host.ingress.tls.port")
		if err != nil {
			log.ErrorF("Failed to get listenSSLPort: %s", err.Error())
			return err
		}

		certFile, err := s.getConfigService().GetAsString("host.ingress.tls.cert")
		if err != nil {
			log.ErrorF("Failed to get certFile: %s", err.Error())
			return err
		}
		paths := s.getPathsService()
		runDir := s.GetKernel().GetRunDirectory()
		if !strings.HasPrefix(certFile, "/") {
			certFile = paths.CominePath(runDir, certFile)
		}

		keyFile, err := s.getConfigService().GetAsString("host.ingress.tls.key")
		if err != nil {
			log.ErrorF("Failed to get keyFile: %s", err.Error())
			return err
		}
		if !strings.HasPrefix(keyFile, "/") {
			keyFile = paths.CominePath(runDir, keyFile)
		}

		go func() {
			// Create a new HTTP server that redirects to HTTPS
			redirectServer := &http.Server{
				Addr: fmt.Sprintf("%s:%d", hostAddr, hostPort),
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					target := fmt.Sprintf("%s%s", rootUrl, r.URL.RequestURI())
					http.Redirect(w, r, target, http.StatusMovedPermanently)
				}),
			}

			log.DebugF("Redirecting HTTP to HTTPS on %s", fmt.Sprintf("%s:%d", hostAddr, hostPort))
			err := redirectServer.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.PanicF("Failed to start redirect server: %s", err.Error())
			}
		}()

		go func() {
			log.DebugF("Listening on %s", rootUrl)
			err = server.RunTLS(fmt.Sprintf("%s:%d", hostAddr, listenSSLPort), certFile, keyFile)
			if err != nil && err != http.ErrServerClosed {
				log.PanicF("Failed to start server: %s", err.Error())
			}
		}()

	} else {
		//HTTP Only listen
		go func() {
			log.DebugF("Listening on %s", rootUrl)
			err := server.Run(fmt.Sprintf("%s:%d", hostAddr, hostPort))

			if err != nil && err != http.ErrServerClosed {
				log.PanicF("Failed to start server: %s", err.Error())
			}
		}()
	}

	return nil

}

func (s *Service) Stop() error {

	log := s.GetLogger()
	server := s.GetServer()
	if server == nil {
		log.Panic("No server found")
	}

	log.Info("Server stopped")
	return nil
}
