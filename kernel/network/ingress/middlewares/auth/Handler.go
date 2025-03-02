package auth

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/config"
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

		accessLevel := accessLevelEvaluation(service, c.Request, default_access, log)
		switch accessLevel {
		case "public":
			if accessLevelPublic(c.Request, log) {
				c.Next()
				return
			}
		case "authenticated":
			if accessLevelAuthenticated(c.Request, log, verify_url, redirect_param) {
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

		default:
			log.ErrorF("Unknown access level: %s", accessLevel)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

	})

	return nil
}

func accessLevelEvaluation(service types.IKernelService, request *http.Request, default_access string, log types.ILogger) string {
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
			return string(policy)
		}
	}

	return default_access
}

func accessLevelPublic(_ *http.Request, _ types.ILogger) bool {
	return true
}

func accessLevelAuthenticated(req *http.Request, log types.ILogger, verifyUrl string, redirectParam string) bool {

	requestUrl := req.URL.Path
	requestQuery := req.URL.RawQuery

	urlVerifyUrl, err := url.Parse(verifyUrl)
	if err != nil {
		log.ErrorF("Failed to parse verifyUrl: %s", err.Error())
		return false
	}

	urlVerifyQuery := urlVerifyUrl.Query()
	urlVerifyQuery.Add(redirectParam, requestUrl)
	urlVerifyQuery.Add("query", requestQuery)
	urlVerifyUrl.RawQuery = urlVerifyQuery.Encode()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 5 * time.Second,
	}

	verifyReq, err := http.NewRequest("GET", urlVerifyUrl.String(), nil)
	if err != nil {
		log.ErrorF("Failed to create request: %s", err.Error())
		return false
	}
	for key, value := range req.Header {
		verifyReq.Header.Add(key, value[0])
	}

	verifyResp, err := client.Do(verifyReq)
	if err != nil {
		log.ErrorF("Failed to execute request: %s", err.Error())
		return false
	}
	if verifyResp.StatusCode == http.StatusOK {
		return true
	}

	return false
}
