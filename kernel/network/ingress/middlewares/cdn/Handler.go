package cdn

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/types"
)

func Register(service types.IKernelService, router *gin.RouterGroup) error {

	log := service.GetLogger().CreateSubLogger("cdn")

	router.GET("/cdn/:package/*file", func(c *gin.Context) {
		pkg := c.Param("package")
		file := c.Param("file")

		unpkg := "https://unpkg.com/" + pkg + file

		body, err := getOrDownloadFile(unpkg)
		if err != nil {
			log.ErrorF("Error downloading file : %s", err.Error())
			c.String(500, "Error downloading file")
			return
		}
		c.Data(200, getContentType(unpkg), body)
	})

	return nil
}

var cache = make(map[string][]byte)

func getOrDownloadFile(url string) ([]byte, error) {

	if _, ok := cache[url]; ok {
		return cache[url], nil
	}

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	cache[url] = body

	return body, nil

}

func getContentType(url string) string {
	ext := url[strings.LastIndex(url, ".")+1:]
	switch ext {
	case "js":
		return "application/javascript"
	case "css":
		return "text/css"
	case "html":
		return "text/html"
	case "png":
		return "image/png"
	case "jpg":
		return "image/jpeg"
	case "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "svg":
		return "image/svg+xml"
	case "json":
		return "application/json"
	case "yaml":
		return "application/yaml"
	}
	return "text/plain"
}
