package plugins

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/tuxounet/k2-sdk/types"
)

const (
	PluginSymbolName = "NewComponent"
	PluginExtension  = ".so"
	PluginEmbedDir   = "dist"
)

func (s *PluginsService) Init() error {
	s.GetLogger().TraceF("[INIT]")

	// Get the embedded filesystem containing plugins
	externals := s.getAppExternals()
	if externals == nil {
		s.GetLogger().TraceF("[INIT] no externals configured, skipping plugin loading")
		return nil
	}

	// Extract plugins from embedded filesystem to lib directory
	err := s.extractPluginsFromEmbed(externals)
	if err != nil {
		return fmt.Errorf("failed to extract plugins: %w", err)
	}

	// Get extracted plugin paths
	pluginPaths, ok := s.GetData("extractedPluginPaths").([]string)
	if !ok || len(pluginPaths) == 0 {
		s.GetLogger().TraceF("[INIT] no plugins to register")
		return nil
	}

	// Load all plugins and register their components
	loadedPlugins := make([]types.PluginInfo, 0)
	for _, pluginPath := range pluginPaths {
		pluginInfo, err := s.loadPlugin(pluginPath)
		if err != nil {
			s.GetLogger().WarnF("[INIT] failed to load plugin %s: %s", filepath.Base(pluginPath), err.Error())
			continue
		}

		loadedPlugins = append(loadedPlugins, *pluginInfo)
		s.GetLogger().InfoF("[INIT] loaded plugin: %s", pluginInfo.Name)
	}

	s.SetData("loadedPlugins", loadedPlugins)
	s.GetLogger().TraceF("[INIT] registered %d plugins", len(loadedPlugins))

	return nil
}

// extractPluginsFromEmbed scans the embedded filesystem and extracts all .so files to lib directory
func (s *PluginsService) extractPluginsFromEmbed(embedFS *embed.FS) error {
	s.GetLogger().TraceF("[EXTRACT] scanning embedded filesystem for plugins")

	// Read the plugin directory from embedded filesystem
	entries, err := fs.ReadDir(embedFS, PluginEmbedDir)
	if err != nil {
		s.GetLogger().TraceF("[EXTRACT] no %s directory in embedded filesystem, skipping: %s", PluginEmbedDir, err.Error())
		return nil
	}

	// Filter .so files
	var pluginFiles []fs.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), PluginExtension) {
			pluginFiles = append(pluginFiles, entry)
		}
	}

	if len(pluginFiles) == 0 {
		s.GetLogger().TraceF("[EXTRACT] no plugins found in embedded filesystem")
		return nil
	}

	// Create lib directory in kernel run directory to extract plugins
	runDir := s.GetKernel().GetRunDirectory()
	libDir := filepath.Join(runDir, "var", "lib")
	if err := os.MkdirAll(libDir, 0755); err != nil {
		return fmt.Errorf("failed to create lib directory: %w", err)
	}
	s.SetData("pluginLibDir", libDir)
	s.GetLogger().TraceF("[EXTRACT] using lib directory: %s", libDir)

	extractedPaths := make([]string, 0)

	for _, entry := range pluginFiles {
		pluginPath, err := s.extractPlugin(embedFS, entry.Name(), libDir)
		if err != nil {
			s.GetLogger().WarnF("[EXTRACT] failed to extract plugin %s: %s", entry.Name(), err.Error())
			continue
		}

		extractedPaths = append(extractedPaths, pluginPath)
		s.GetLogger().TraceF("[EXTRACT] extracted plugin: %s", entry.Name())
	}

	s.GetLogger().TraceF("[EXTRACT] extracted %d plugins", len(extractedPaths))
	s.SetData("extractedPluginPaths", extractedPaths)
	return nil
}

// extractPlugin extracts a plugin from the embedded filesystem to lib directory
func (s *PluginsService) extractPlugin(embedFS *embed.FS, pluginName string, libDir string) (string, error) {
	embedPath := filepath.Join(PluginEmbedDir, pluginName)

	// Read plugin content from embedded filesystem
	content, err := fs.ReadFile(embedFS, embedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded plugin: %w", err)
	}

	// Write plugin to lib directory
	pluginPath := filepath.Join(libDir, pluginName)
	err = os.WriteFile(pluginPath, content, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to write plugin to lib directory: %w", err)
	}

	return pluginPath, nil
}

// loadPlugin loads a single plugin .so file and registers its component
func (s *PluginsService) loadPlugin(pluginPath string) (*types.PluginInfo, error) {
	s.GetLogger().TraceF("[LOAD] loading plugin from: %s", pluginPath)

	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		s.GetLogger().ErrorF("[LOAD] failed to open plugin %s: %s", pluginPath, err.Error())
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}
	s.GetLogger().DebugF("[LOAD] successfully opened plugin: %s", pluginPath)

	// Lookup the NewComponent symbol
	sym, err := p.Lookup(PluginSymbolName)
	if err != nil {
		s.GetLogger().ErrorF("[LOAD] plugin %s does not export %s symbol: %s", pluginPath, PluginSymbolName, err.Error())
		return nil, fmt.Errorf("plugin does not export %s symbol: %w", PluginSymbolName, err)
	}
	s.GetLogger().DebugF("[LOAD] found symbol %s in plugin", PluginSymbolName)

	// Assert the symbol to the expected function signature
	componentCtor, ok := sym.(func(types.IApp) types.IAppComponent)
	if !ok {
		s.GetLogger().ErrorF("[LOAD] plugin %s has invalid signature, expected func(types.IApp) types.IAppComponent", pluginPath)
		return nil, fmt.Errorf("plugin %s has invalid signature, expected func(types.IApp) types.IAppComponent", PluginSymbolName)
	}
	s.GetLogger().DebugF("[LOAD] validated symbol signature for plugin: %s", pluginPath)

	// Register the component constructor with the app
	app := s.GetKernel().GetApp()
	app.AddComponent(componentCtor)
	s.GetLogger().DebugF("[LOAD] registered component constructor with app")

	// Extract plugin name from filename
	baseName := filepath.Base(pluginPath)
	name := baseName[:len(baseName)-len(PluginExtension)]
	s.GetLogger().InfoF("[LOAD] successfully loaded plugin: %s", name)

	return &types.PluginInfo{
		Name: name,
		Path: pluginPath,
	}, nil
}
