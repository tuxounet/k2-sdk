package host

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/tuxounet/k2-sdk/types"
)

const (
	PluginSymbolName = "NewComponent"
	PluginExtension  = ".so"
)

func (s *PluginsHostService) Init() error {
	s.GetLogger().TraceF("[INIT]")

	// Discover and load plugins
	pluginDir := s.getPluginDirectory()
	if pluginDir == "" {
		s.GetLogger().TraceF("[INIT] no plugin directory configured, skipping plugin loading")
		return nil
	}

	err := s.discoverAndLoadPlugins(pluginDir)
	if err != nil {
		return fmt.Errorf("failed to discover and load plugins: %w", err)
	}

	s.GetLogger().TraceF("[INIT] complete")
	return nil
}

func (s *PluginsHostService) Register() error {
	s.GetLogger().TraceF("[REGISTER]")
	return nil
}

func (s *PluginsHostService) Start() error {
	s.GetLogger().TraceF("[START]")
	return nil
}

func (s *PluginsHostService) Stop() error {
	s.GetLogger().TraceF("[STOP]")
	return nil
}

// getPluginDirectory returns the configured plugin directory
func (s *PluginsHostService) getPluginDirectory() string {
	// Check environment variable first
	if dir := os.Getenv("K2_PLUGIN_DIR"); dir != "" {
		return dir
	}

	// Check service config
	if dir := s.GetConfig("plugin_dir"); dir != "" {
		return dir
	}

	// Default to plugins subdirectory in run directory
	runDir := s.GetKernel().GetRunDirectory()
	defaultPluginDir := filepath.Join(runDir, "plugins")

	// Only return if directory exists
	if _, err := os.Stat(defaultPluginDir); err == nil {
		return defaultPluginDir
	}

	return ""
}

// discoverAndLoadPlugins scans the plugin directory and loads all .so files
func (s *PluginsHostService) discoverAndLoadPlugins(pluginDir string) error {
	s.GetLogger().TraceF("[DISCOVER] scanning plugin directory: %s", pluginDir)

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			s.GetLogger().TraceF("[DISCOVER] plugin directory does not exist, skipping")
			return nil
		}
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	loadedPlugins := make([]types.PluginInfo, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != PluginExtension {
			continue
		}

		pluginPath := filepath.Join(pluginDir, entry.Name())
		pluginInfo, err := s.loadPlugin(pluginPath)
		if err != nil {
			s.GetLogger().WarnF("[DISCOVER] failed to load plugin %s: %s", entry.Name(), err.Error())
			continue
		}

		loadedPlugins = append(loadedPlugins, *pluginInfo)
		s.GetLogger().InfoF("[DISCOVER] loaded plugin: %s", pluginInfo.Name)
	}

	s.SetData("loadedPlugins", loadedPlugins)
	s.GetLogger().TraceF("[DISCOVER] loaded %d plugins", len(loadedPlugins))

	return nil
}

// loadPlugin loads a single plugin .so file and registers its component
func (s *PluginsHostService) loadPlugin(pluginPath string) (*types.PluginInfo, error) {
	s.GetLogger().TraceF("[LOAD] loading plugin from: %s", pluginPath)

	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Lookup the NewComponent symbol
	sym, err := p.Lookup(PluginSymbolName)
	if err != nil {
		return nil, fmt.Errorf("plugin does not export %s symbol: %w", PluginSymbolName, err)
	}

	// Assert the symbol to the expected function signature
	componentCtor, ok := sym.(func(types.IApp) types.IAppComponent)
	if !ok {
		return nil, fmt.Errorf("plugin %s has invalid signature, expected func(types.IApp) types.IAppComponent", PluginSymbolName)
	}

	// Register the component constructor with the app
	app := s.GetKernel().GetApp()
	app.AddComponent(componentCtor)

	// Extract plugin name from filename
	pluginName := filepath.Base(pluginPath)
	pluginName = pluginName[:len(pluginName)-len(PluginExtension)]

	return &types.PluginInfo{
		Name: pluginName,
		Path: pluginPath,
	}, nil
}

// GetLoadedPlugins returns the list of loaded plugins
func (s *PluginsHostService) GetLoadedPlugins() []types.PluginInfo {
	data := s.GetData("loadedPlugins")
	if data == nil {
		return []types.PluginInfo{}
	}
	return data.([]types.PluginInfo)
}
