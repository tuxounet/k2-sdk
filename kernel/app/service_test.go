package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/testutils"
	"github.com/tuxounet/k2-sdk/types"
)

func newTestAppService(t *testing.T) (*Service, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "k2-app-test-*")
	require.NoError(t, err)

	kernel := testutils.NewMockKernel("test-app", "1.0.0")
	kernel.SetRunDir(tmpDir)

	// Create config service (dependency)
	configSvc := config.NewService(kernel).(*config.Service)
	kernel.SetService(configSvc)
	configSvc.SetData("records", make(map[string]any))

	// Create mock app with components
	app := testutils.NewMockApp("test-app", "1.0.0")
	comp := testutils.NewMockComponent(app, "test-comp", 0)
	app.SetComponents([]types.IAppComponent{comp})
	kernel.SetApp(app)

	appSvc := NewService(kernel).(*Service)
	return appSvc, tmpDir
}

func TestAppService_NewService(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.Equal(t, ServiceKey, svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
}

func TestAppService_Init(t *testing.T) {
	svc, tmpDir := newTestAppService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	// After init, the app should have kernel and logger set
	app := svc.GetKernel().GetApp()
	assert.NotNil(t, app.GetKernel())
	assert.NotNil(t, app.GetLogger())
}

func TestAppService_Start(t *testing.T) {
	svc, tmpDir := newTestAppService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	err = svc.Start()
	require.NoError(t, err)
}

func TestAppService_Stop(t *testing.T) {
	svc, tmpDir := newTestAppService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	err = svc.Start()
	require.NoError(t, err)

	err = svc.Stop()
	require.NoError(t, err)
}

func TestAppService_Init_LoadsEnvVars(t *testing.T) {
	svc, tmpDir := newTestAppService(t)
	defer os.RemoveAll(tmpDir)

	os.Setenv("K2_TEST_VAR", "test_value")
	defer os.Unsetenv("K2_TEST_VAR")

	err := svc.Init()
	require.NoError(t, err)

	configSvc := svc.getConfigService()
	assert.Equal(t, "test_value", configSvc.Get("test.var"))
}
