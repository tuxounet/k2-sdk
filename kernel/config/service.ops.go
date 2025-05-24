package config

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"dario.cat/mergo"
	"github.com/tuxounet/k2-sdk/system"
)

//go:embed defaults/*
var defaultConfig embed.FS

func (s *Service) initDefaultConfig() error {

	s.SetData("records", make(map[string]any))
	err := s.LoadFromEmbedFS("kernel", "defaults", &defaultConfig)
	if err != nil {
		s.GetLogger().ErrorF("Failed to load default configuration: %v", err)
		return err
	}

	return nil
}

func (s *Service) LoadFromEmbedFS(source string, folder string, fs *embed.FS) error {

	if fs == nil {
		//nothing to do
		return nil
	}

	entries, err := fs.ReadDir(folder)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // directory does not exist, nothing to load
		} else {
			return fmt.Errorf("failed to read embeded config directory %s: %v", folder, err)
		}
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(folder + "/" + entry.Name())
		if err != nil {
			s.GetLogger().ErrorF("Failed to read config file %v: %v", entry.Name(), err)
			return err
		}
		configKey := fmt.Sprintf("%s/%s/%s", source, folder, entry.Name())

		//unmarshal data
		result, err := system.LoadYamlFromString[map[string]any](string(data))
		if err != nil {
			s.GetLogger().ErrorF("Failed to unmarshal config file %s: %s", configKey, err.Error())
			return err
		}

		current := s.GetCurrent()
		err = mergo.Merge(&current, result, mergo.WithOverride)
		if err != nil {
			s.GetLogger().ErrorF("Failed to merge configuration with %s: %s", configKey, err.Error())
			return err
		}
		s.GetLogger().DebugF("Loaded config file %s", configKey)
		s.SetData("records", current)

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
