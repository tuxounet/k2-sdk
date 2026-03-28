package bases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuxounet/k2-sdk/testutils"
)

func TestNewBaseKernelService(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "my-service")

	assert.Equal(t, "my-service", svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
	assert.NotNil(t, svc.GetLogger())
}

func TestBaseKernelService_Config_GetSet(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "svc")

	svc.SetConfig("host", "localhost")
	svc.SetConfig("port", "8080")

	assert.Equal(t, "localhost", svc.GetConfig("host"))
	assert.Equal(t, "8080", svc.GetConfig("port"))
}

func TestBaseKernelService_Config_MissingKey(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "svc")

	assert.Equal(t, "", svc.GetConfig("nonexistent"))
}

func TestBaseKernelService_Config_Overwrite(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "svc")

	svc.SetConfig("key", "original")
	svc.SetConfig("key", "updated")

	assert.Equal(t, "updated", svc.GetConfig("key"))
}

func TestBaseKernelService_Data_GetSet(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "svc")

	svc.SetData("items", []string{"a", "b"})
	svc.SetData("count", 42)

	assert.Equal(t, []string{"a", "b"}, svc.GetData("items"))
	assert.Equal(t, 42, svc.GetData("count"))
}

func TestBaseKernelService_Data_MissingKey(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "svc")

	assert.Nil(t, svc.GetData("nonexistent"))
}

func TestBaseKernelService_Lifecycle(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "svc")

	assert.NoError(t, svc.Init())
	assert.NoError(t, svc.Register())
	assert.NoError(t, svc.Start())
	assert.NoError(t, svc.Stop())
}

func TestBaseKernelService_Logger_SubLogger(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewBaseKernelService(kernel, "my-svc")

	logger := svc.GetLogger()
	assert.Contains(t, logger.GetName(), "my-svc")
}
