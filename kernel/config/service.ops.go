package config

import (
	"embed"
	"strings"

	"dario.cat/mergo"
	"github.com/tuxounet/k2-sdk/system"
)

//go:embed defaults/*
var defaultConfig embed.FS

func (s *Service) initDefaultConfig() error {

	s.SetData("records", make(map[string]interface{}))
	err := s.LoadFromEmbedFS("defaults", &defaultConfig)
	if err != nil {
		s.GetLogger().ErrorF("Failed to load default configuration: %v", err)
		return err
	}

	return nil
}

func (s *Service) LoadFromEmbedFS(folder string, fs *embed.FS) error {

	if fs != nil {

		entries, err := fs.ReadDir(folder)
		if err != nil {
			s.GetLogger().WarnF("Failed to read config directory: %v", err)
		} else {

			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
					continue
				}
				data, err := fs.ReadFile(folder + "/" + entry.Name())
				if err != nil {
					s.GetLogger().ErrorF("Failed to read config file %v: %v", entry.Name(), err)
					return err
				}

				//unmarshal data
				result, err := system.LoadYamlFromString[map[string]interface{}](string(data))
				if err != nil {
					s.GetLogger().ErrorF("Failed to unmarshal config file %v: %v", entry.Name(), err)
					return err
				}

				current := s.GetCurrent()
				err = mergo.Merge(&current, result, mergo.WithOverride)
				if err != nil {
					s.GetLogger().ErrorF("Failed to merge configuration: %v", err)
					return err
				}
				s.SetData("records", current)

			}

		}
	}
	return nil
}

type untemplateData struct {
	Config     map[string]interface{}
	FileName   string
	Name       string
	Version    string
	RootUrl    string
	BasePath   string
	UIBasePath string
}

func (s *Service) Untemplate(templated []byte, baseRoute string) ([]byte, error) {
	app := s.GetKernel().GetApp()

	rootUrl, err := s.GetAsString("host.ingress.rootUrl")
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(baseRoute, "/") {
		baseRoute += "/"
	}
	uiPathSuffix := "ui/"
	uiBaseRoute := baseRoute + uiPathSuffix

	templateValues := untemplateData{
		Config:     s.GetCurrent(),
		Name:       app.GetName(),
		Version:    app.GetVersion(),
		RootUrl:    rootUrl,
		BasePath:   baseRoute,
		UIBasePath: uiBaseRoute,
	}
	result, err := system.UnTemplateWithGoTemplate(string(templated), templateValues)

	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}
