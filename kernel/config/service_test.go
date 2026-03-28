package config

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newTestConfigService() *Service {
	kernel := testutils.NewMockKernel("test-app", "1.0")
	svc := NewService(kernel).(*Service)
	svc.SetData("records", make(map[string]any))
	return svc
}

func newTestConfigServiceWithData(data map[string]any) *Service {
	svc := newTestConfigService()
	svc.SetData("records", data)
	return svc
}

// --- Get ---

func TestConfig_Get_SimpleKey(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{
		"name": "test",
	})
	assert.Equal(t, "test", svc.Get("name"))
}

func TestConfig_Get_DotNotation(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{
		"host": map[string]any{
			"ingress": map[string]any{
				"port": 8080,
			},
		},
	})
	assert.Equal(t, 8080, svc.Get("host.ingress.port"))
}

func TestConfig_Get_MissingKey(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{
		"name": "test",
	})
	assert.Nil(t, svc.Get("nonexistent"))
}

func TestConfig_Get_MissingNestedKey(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{
		"host": map[string]any{
			"port": 8080,
		},
	})
	assert.Nil(t, svc.Get("host.missing.deep"))
}

func TestConfig_Get_CaseInsensitive(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{
		"host": map[string]any{
			"port": 8080,
		},
	})
	assert.Equal(t, 8080, svc.Get("HOST.PORT"))
	assert.Equal(t, 8080, svc.Get("Host.Port"))
}

func TestConfig_Get_NonMapValue(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{
		"name": "test",
	})
	assert.Nil(t, svc.Get("name.sub"))
}

// --- Has ---

func TestConfig_Has_Existing(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{"key": "val"})
	assert.True(t, svc.Has("key"))
}

func TestConfig_Has_Missing(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{})
	assert.False(t, svc.Has("nonexistent"))
}

// --- GetAsString ---

func TestConfig_GetAsString_Success(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{"name": "hello"})

	val, err := svc.GetAsString("name")
	require.NoError(t, err)
	assert.Equal(t, "hello", val)
}

func TestConfig_GetAsString_Missing(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{})

	_, err := svc.GetAsString("missing")
	assert.Error(t, err)
}

// --- GetAsInt ---

func TestConfig_GetAsInt_Success(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{"port": 8080})

	val, err := svc.GetAsInt("port")
	require.NoError(t, err)
	assert.Equal(t, 8080, val)
}

func TestConfig_GetAsInt_Missing(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{})

	val, err := svc.GetAsInt("missing")
	assert.Error(t, err)
	assert.Equal(t, -1, val)
}

// --- GetAsBool ---

func TestConfig_GetAsBool_Success(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{"enabled": true})

	val, err := svc.GetAsBool("enabled")
	require.NoError(t, err)
	assert.True(t, val)
}

func TestConfig_GetAsBool_Missing(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{})

	val, err := svc.GetAsBool("missing")
	assert.Error(t, err)
	assert.False(t, val)
}

// --- GetAsStringOrDefault ---

func TestConfig_GetAsStringOrDefault_Found(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{"name": "found"})

	val, err := svc.GetAsStringOrDefault("name", "default")
	require.NoError(t, err)
	assert.Equal(t, "found", val)
}

func TestConfig_GetAsStringOrDefault_Default(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{})

	val, err := svc.GetAsStringOrDefault("missing", "fallback")
	require.NoError(t, err)
	assert.Equal(t, "fallback", val)
}

// --- GetAsIntOrDefault ---

func TestConfig_GetAsIntOrDefault_Found(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{"port": 9090})

	val, err := svc.GetAsIntOrDefault("port", 8080)
	require.NoError(t, err)
	assert.Equal(t, 9090, val)
}

func TestConfig_GetAsIntOrDefault_Default(t *testing.T) {
	svc := newTestConfigServiceWithData(map[string]any{})

	val, err := svc.GetAsIntOrDefault("missing", 8080)
	require.NoError(t, err)
	assert.Equal(t, 8080, val)
}

// --- SetValue ---

func TestConfig_SetValue_Simple(t *testing.T) {
	svc := newTestConfigService()

	err := svc.SetValue("name", "test")
	require.NoError(t, err)
	assert.Equal(t, "test", svc.Get("name"))
}

func TestConfig_SetValue_Nested(t *testing.T) {
	svc := newTestConfigService()

	err := svc.SetValue("host.port", 8080)
	require.NoError(t, err)
	assert.Equal(t, 8080, svc.Get("host.port"))
}

func TestConfig_SetValue_DeepNested(t *testing.T) {
	svc := newTestConfigService()

	err := svc.SetValue("a.b.c.d", "deep")
	require.NoError(t, err)
	assert.Equal(t, "deep", svc.Get("a.b.c.d"))
}

func TestConfig_SetValue_EmptyKey(t *testing.T) {
	svc := newTestConfigService()

	err := svc.SetValue("", "value")
	assert.Error(t, err)
}

func TestConfig_SetValue_Overwrite(t *testing.T) {
	svc := newTestConfigService()

	svc.SetValue("key", "original")
	svc.SetValue("key", "updated")

	assert.Equal(t, "updated", svc.Get("key"))
}

func TestConfig_SetValue_AutoVivification(t *testing.T) {
	svc := newTestConfigService()

	err := svc.SetValue("new.nested.key", "value")
	require.NoError(t, err)

	assert.Equal(t, "value", svc.Get("new.nested.key"))
	assert.NotNil(t, svc.Get("new"))
	assert.NotNil(t, svc.Get("new.nested"))
}

// --- GetCurrent ---

func TestConfig_GetCurrent(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2}
	svc := newTestConfigServiceWithData(data)

	current := svc.GetCurrent()
	assert.Equal(t, 1, current["a"])
	assert.Equal(t, 2, current["b"])
}

// --- LoadFromEnvVars ---

func TestConfig_LoadFromEnvVars(t *testing.T) {
	svc := newTestConfigService()

	os.Setenv("K2_HOST_PORT", "9090")
	os.Setenv("K2_APP_NAME", "myapp")
	defer os.Unsetenv("K2_HOST_PORT")
	defer os.Unsetenv("K2_APP_NAME")

	err := svc.LoadFromEnvVars("test")
	require.NoError(t, err)

	assert.Equal(t, "9090", svc.Get("host.port"))
	assert.Equal(t, "myapp", svc.Get("app.name"))
}

func TestConfig_LoadFromEnvVars_NoK2Vars(t *testing.T) {
	svc := newTestConfigService()

	err := svc.LoadFromEnvVars("test")
	require.NoError(t, err)
}

// --- LoadFromEmbedFS ---

func TestConfig_LoadFromEmbedFS_NilFS(t *testing.T) {
	svc := newTestConfigService()

	err := svc.LoadFromEmbedFS("test", "config", nil)
	require.NoError(t, err)
}

// We cannot test with a real embed.FS from a test, but we can test the default config init
func TestConfig_InitDefaultConfig(t *testing.T) {
	svc := newTestConfigService()

	err := svc.initDefaultConfig()
	require.NoError(t, err)

	current := svc.GetCurrent()
	assert.NotNil(t, current)
}

// --- LoadFromEmbedFS with fstest (testing/fstest doesn't implement embed.FS) ---
// embed.FS is a concrete type, not an interface, so we test via initDefaultConfig
// which uses the embedded defaults

func TestConfig_DefaultConfig_HasHostSection(t *testing.T) {
	svc := newTestConfigService()

	err := svc.initDefaultConfig()
	require.NoError(t, err)

	// The default config should have a "host" section
	assert.NotNil(t, svc.Get("host"))
}

// Note: fstest.MapFS is available but embed.FS is a concrete type.
// This variable is referenced to satisfy the import checker.
var _ = fstest.MapFS{}
