package levels

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/tuxounet/k2-sdk/types"
)

func AllowAccessLevelAuthenticated(req *http.Request, log types.ILogger, verifyUrl string, redirectParam string) bool {

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

func RedirectAccessLevelAuthenticatedLogin(req *http.Request, log types.ILogger, loginUrl string, redirectParam string) bool {

	requestUrl := req.URL.Path
	requestQuery := req.URL.RawQuery

	urlLoginUrl, err := url.Parse(loginUrl)
	if err != nil {
		log.ErrorF("Failed to parse loginUrl: %s", err.Error())
		return false
	}

	urlLoginQuery := urlLoginUrl.Query()
	urlLoginQuery.Add(redirectParam, requestUrl)
	urlLoginQuery.Add("query", requestQuery)
	urlLoginUrl.RawQuery = urlLoginQuery.Encode()

	log.InfoF("Redirecting to login: %s", urlLoginUrl.String())
	req.URL.Path = urlLoginUrl.String()
	return true
}
