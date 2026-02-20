package types

// PluginComponentCtor is the expected function signature exported by plugins
// Plugins must export a function named "NewComponent" with this signature
type PluginComponentCtor = func(app IApp) IAppComponent

// PluginInfo contains metadata about a loaded plugin
type PluginInfo struct {
	Name string
	Path string
}
