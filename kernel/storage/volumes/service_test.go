package volumes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuxounet/k2-sdk/testutils"
)

func TestVolumesService_NewService(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.Equal(t, ServiceKey, svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
}

func TestVolumesService_Lifecycle(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.NoError(t, svc.Init())
	assert.NoError(t, svc.Register())
	assert.NoError(t, svc.Start())
	assert.NoError(t, svc.Stop())
}
