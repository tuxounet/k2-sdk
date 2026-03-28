package logging

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
)

func TestNewService(t *testing.T) {
	kernel := testutils.NewMockKernel("test-app", "1.0")
	svc := NewService(kernel)

	assert.Equal(t, LoggingServiceKey, svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
}

func TestLoggingService_Init(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "k2-logging-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	kernel := testutils.NewMockKernel("test-app", "1.0")
	kernel.SetRunDir(tmpDir)

	app := testutils.NewMockApp("test-app", "1.0")
	kernel.SetApp(app)

	svc := NewService(kernel)

	err = svc.Init()
	require.NoError(t, err)

	rootLogger := svc.GetRootLogger()
	assert.NotNil(t, rootLogger)

	serviceLogger := svc.GetLogger()
	assert.NotNil(t, serviceLogger)
}

func TestLoggingService_Init_Idempotent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "k2-logging-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	kernel := testutils.NewMockKernel("test-app", "1.0")
	kernel.SetRunDir(tmpDir)

	app := testutils.NewMockApp("test-app", "1.0")
	kernel.SetApp(app)

	svc := NewService(kernel)

	err = svc.Init()
	require.NoError(t, err)

	logger1 := svc.GetRootLogger()

	err = svc.Init()
	require.NoError(t, err)

	logger2 := svc.GetRootLogger()
	assert.Equal(t, logger1, logger2) // Same instance, not recreated
}

func TestLoggingService_ConfigGetSet(t *testing.T) {
	kernel := testutils.NewMockKernel("test-app", "1.0")
	svc := NewService(kernel)

	svc.SetConfig("key", "value")
	assert.Equal(t, "value", svc.GetConfig("key"))
	assert.Equal(t, "", svc.GetConfig("missing"))
}

func TestLoggingService_DataGetSet(t *testing.T) {
	kernel := testutils.NewMockKernel("test-app", "1.0")
	svc := NewService(kernel)

	svc.SetData("items", []int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, svc.GetData("items"))
	assert.Nil(t, svc.GetData("missing"))
}

func TestLoggingService_Lifecycle(t *testing.T) {
	kernel := testutils.NewMockKernel("test-app", "1.0")
	svc := NewService(kernel)

	assert.NoError(t, svc.Register())
	assert.NoError(t, svc.Start())
	assert.NoError(t, svc.Stop())
}
