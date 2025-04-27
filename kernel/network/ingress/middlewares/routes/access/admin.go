package access

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/types"
)

func isAdminEnabled(configService *config.Service) bool {
	admin_enabled, err := configService.GetAsBool("host.ingress.auth.admin.enabled")
	if err != nil {
		return false
	}
	return admin_enabled
}

func AllowAccessLevelAdmin(req *http.Request, log types.ILogger, configService *config.Service) bool {
	requestUrl := req.URL.Path
	requestQuery := req.URL.RawQuery

	if !isAdminEnabled(configService) {
		log.InfoF("Admin access is disabled for path: %s", requestUrl)
		return false
	}
	verify_url, err := configService.GetAsString("host.ingress.auth.admin.verifyUrl")
	if err != nil {
		log.ErrorF("Failed to get verifyUrl: %s", err.Error())
		return false
	}
	redirect_param, err := configService.GetAsString("host.ingress.auth.admin.redirectParam")
	if err != nil {
		log.ErrorF("Failed to get redirectParam: %s", err.Error())
		return false
	}

	if !strings.HasPrefix(verify_url, "http") {
		rootUrl, err := configService.GetAsString("host.ingress.rootUrl")
		if err != nil {
			log.ErrorF("Failed to get rootUrl: %s", err.Error())
			return false
		}
		verify_url = fmt.Sprintf("%s%s", rootUrl, verify_url)

	}

	urlVerifyUrl, err := url.Parse(verify_url)
	if err != nil {
		log.ErrorF("Failed to parse verifyUrl: %s", err.Error())
		return false
	}

	urlVerifyQuery := urlVerifyUrl.Query()
	urlVerifyQuery.Add(redirect_param, requestUrl)
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

func RedirectAccessLevelAdminLogin(req *http.Request, log types.ILogger, configService *config.Service, ctx *gin.Context) {
	requestUrl := req.URL.Path
	requestQuery := req.URL.RawQuery

	if !isAdminEnabled(configService) {
		log.InfoF("Admin access is disabled for path: %s", requestUrl)
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	login_url, err := configService.GetAsString("host.ingress.auth.admin.loginUrl")
	if err != nil {
		log.ErrorF("Failed to get loginUrl: %s", err.Error())
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	loginUrl, err := url.Parse(login_url)
	if err != nil {
		log.ErrorF("Failed to parse loginUrl: %s", err.Error())
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	redirect_param, err := configService.GetAsString("host.ingress.auth.admin.redirectParam")
	if err != nil {
		log.ErrorF("Failed to get redirectParam: %s", err.Error())
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	urlLoginQuery := loginUrl.Query()
	urlLoginQuery.Add(redirect_param, requestUrl)
	urlLoginQuery.Add("query", requestQuery)
	loginUrl.RawQuery = urlLoginQuery.Encode()

	log.InfoF("Redirecting to login: %s", loginUrl.String())
	req.URL.Path = loginUrl.String()
	ctx.Redirect(http.StatusFound, loginUrl.String())
	ctx.Abort()
}
