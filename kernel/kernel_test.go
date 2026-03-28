package kernel

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
	"github.com/tuxounet/k2-sdk/types"
)

func TestNewKernelRuntime(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)
	require.NotNil(t, kernel)

	assert.Equal(t, "test-app", kernel.GetName())
	assert.Equal(t, "1.0.0", kernel.GetVersion())
	assert.Equal(t, tmpDir, kernel.GetRunDirectory())
	assert.False(t, kernel.IsUnsecure())
	assert.NotNil(t, kernel.GetLogger())
	assert.NotNil(t, kernel.GetRootContext())
	assert.Equal(t, app, kernel.GetApp())
}

func TestNewKernelRuntime_Unsecure(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", true)
	assert.True(t, kernel.IsUnsecure())
}

func TestKernelRuntime_ServiceRegistry(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)

	// Verify core services are registered
	assert.NotNil(t, kernel.GetService("monitoring.logging"))
	assert.NotNil(t, kernel.GetService("storage.paths"))
	assert.NotNil(t, kernel.GetService("config"))
	assert.NotNil(t, kernel.GetService("secrets"))
	assert.NotNil(t, kernel.GetService("scheduler"))
	assert.NotNil(t, kernel.GetService("network.ingress.http"))
}

func TestKernelRuntime_GetService_Missing(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)

	svc := kernel.GetService("nonexistent.service")
	assert.Nil(t, svc)
}

func TestKernelRuntime_DefaultRunDir(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)
	require.NotNil(t, kernel)

	// Should default to cwd/.data
	assert.Contains(t, kernel.GetRunDirectory(), ".data")

	// Cleanup
	os.RemoveAll(kernel.GetRunDirectory())
}

func TestKernelRuntime_AllServices(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)

	// Verify all 12 services are registered
	expectedServices := []string{
		"monitoring.logging",
		"storage.paths",
		"profile",
		"storage.volumes",
		"storage.stores",
		"config",
		"secrets",
		"compute",
		"plugins",
		"app",
		"scheduler",
		"network.ingress.http",
	}

	for _, key := range expectedServices {
		svc := kernel.GetService(types.KernelServiceContextKey(key))
		assert.NotNilf(t, svc, "service %s should be registered", key)
		if svc != nil {
			assert.Equal(t, key, svc.GetName())
		}
	}
}

func TestKernelRuntime_SetService(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)

	mockSvc := testutils.NewMockKernelService(kernel, "custom.test")
	kernel.SetService(mockSvc)

	retrieved := kernel.GetService("custom.test")
	assert.NotNil(t, retrieved)
	assert.Equal(t, "custom.test", retrieved.GetName())
}

func TestKernelRuntime_Init(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)

	err = kernel.Init()
	require.NoError(t, err)
}

func TestKernelRuntime_Init_Register(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)

	err = kernel.Init()
	require.NoError(t, err)

	err = kernel.Register()
	require.NoError(t, err)
}

func TestKernelRuntime_FullLifecycle(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	tmpDir, err := os.MkdirTemp("", "k2-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	os.Setenv("RUN_DIR", tmpDir)
	defer os.Unsetenv("RUN_DIR")

	kernel := NewKernelRuntime(app, "1.0.0", false)

	err = kernel.Init()
	require.NoError(t, err)

	err = kernel.Register()
	require.NoError(t, err)

	err = kernel.Start()
	require.NoError(t, err)

	err = kernel.Stop()
	require.NoError(t, err)
}
