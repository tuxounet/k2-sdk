package files

import (
	"embed"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/types"
)

func ServeAppMiddleware(app types.IApp, log types.ILogger, baseRoute string, fs *embed.FS, fsPrefix string) (gin.HandlerFunc, error) {
	configService := app.GetKernel().GetService(config.ServiceKey).(*config.Service)

	return serveMiddleware(configService, log, baseRoute, fs, fsPrefix)
}

func ServeComponentMiddleware(component types.IAppComponent, log types.ILogger, baseRoute string, fs *embed.FS, fsPrefix string) (gin.HandlerFunc, error) {
	configService := component.GetApp().GetKernel().GetService(config.ServiceKey).(*config.Service)

	return serveMiddleware(configService, log, baseRoute, fs, fsPrefix)
}

func serveMiddleware(configService *config.Service, log types.ILogger, baseRoute string, fs *embed.FS, fsPrefix string) (gin.HandlerFunc, error) {

	fileMap := make(map[string][]byte)

	files, err := WalkFolder(fsPrefix, fs)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filename := path.Base(file)

		fileBody, err := fs.ReadFile(file)
		if err != nil {
			return nil, err
		}

		ext := GetExentionFromFilename(filename)
		switch ext {
		case "html":
			fileBody, err = configService.Untemplate(fileBody, baseRoute)
			if err != nil {
				log.ErrorF("Error untamplating %s: %s", file, err.Error())
				return nil, err
			}
		}
		relativeFile := strings.TrimPrefix(file, fsPrefix)
		fileMap[relativeFile] = fileBody
	}

	return func(c *gin.Context) {
		path := c.Param("any")
		ext := GetExentionFromFilename(path)
		contentType := GetContentTypeFromExtension(ext)
		fileBody, found := fileMap[path]
		if !found {
			//default content
			defaultContentPage := "/index.html"
			fileBody, found = fileMap[defaultContentPage]
			if !found {
				c.Status(404)
				return
			}
			contentType = GetContentTypeFromExtension("html")
			c.Data(200, contentType, fileBody)
			return
		}

		c.Data(200, contentType, fileBody)
	}, nil

}
