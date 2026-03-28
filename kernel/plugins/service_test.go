package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
	"github.com/tuxounet/k2-sdk/types"
)

func TestPluginsService_NewService(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.Equal(t, ServiceKey, svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
}

func TestPluginsService_Init_NoExternals(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	app := testutils.NewMockApp("test", "1.0")
	kernel.SetApp(app)
	// app.GetExternals() returns nil by default

	svc := NewService(kernel).(*PluginsService)
	err := svc.Init()
	require.NoError(t, err)
}

func TestPluginsService_GetLoadedPlugins_Empty(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*PluginsService)

	plugins := svc.GetLoadedPlugins()
	assert.Empty(t, plugins)
}

func TestPluginsService_GetLoadedPlugins_WithData(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*PluginsService)

	svc.SetData("loadedPlugins", []types.PluginInfo{
		{Name: "plugin1", Path: "/tmp/plugin1.so"},
		{Name: "plugin2", Path: "/tmp/plugin2.so"},
	})

	plugins := svc.GetLoadedPlugins()
	assert.Len(t, plugins, 2)
	assert.Equal(t, "plugin1", plugins[0].Name)
	assert.Equal(t, "plugin2", plugins[1].Name)
}

func TestPluginsService_Constants(t *testing.T) {
	assert.Equal(t, "NewComponent", PluginSymbolName)
	assert.Equal(t, ".so", PluginExtension)
	assert.Equal(t, "dist", PluginEmbedDir)
}

func TestPluginsService_Lifecycle(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*PluginsService)

	assert.NoError(t, svc.Register())
	assert.NoError(t, svc.Start())
	assert.NoError(t, svc.Stop())
}
