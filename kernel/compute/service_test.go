package compute

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/compute/types"
	"github.com/tuxounet/k2-sdk/kernel/config"
	"github.com/tuxounet/k2-sdk/kernel/network/ingress"
	"github.com/tuxounet/k2-sdk/kernel/profile"
	"github.com/tuxounet/k2-sdk/kernel/storage/paths"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newTestComputeService(t *testing.T, enabled bool) (*Service, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "k2-compute-test-*")
	require.NoError(t, err)

	kernel := testutils.NewMockKernel("test-app", "1.0.0")
	kernel.SetRunDir(tmpDir)

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	// Config service
	configSvc := config.NewService(kernel).(*config.Service)
	kernel.SetService(configSvc)
	configSvc.SetData("records", map[string]any{
		"host": map[string]any{
			"compute": map[string]any{
				"enabled": enabled,
			},
			"ingress": map[string]any{
				"root":    "http://localhost:8080",
				"address": "0.0.0.0",
				"port":    8080,
			},
		},
	})

	// Paths service
	pathsSvc := paths.NewService(kernel)
	kernel.SetService(pathsSvc)

	// Profile service (dependency for stores)
	profileSvc := profile.NewService(kernel).(*profile.ProfileService)
	kernel.SetService(profileSvc)
	err = profileSvc.Init()
	require.NoError(t, err)

	// Stores service (dependency for compute Register)
	storesSvc := stores.NewService(kernel).(*stores.Service)
	kernel.SetService(storesSvc)
	err = storesSvc.Init()
	require.NoError(t, err)

	// Ingress service
	ingressSvc := ingress.NewService(kernel)
	kernel.SetService(ingressSvc)

	svc := NewService(kernel).(*Service)
	return svc, tmpDir
}

func TestComputeService_NewService(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.Equal(t, ServiceKey, svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
}

func TestComputeService_Init_Disabled(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, false)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)
}

func TestComputeService_Register_Disabled(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, false)
	defer os.RemoveAll(tmpDir)

	// Register should work even when disabled (just renders inventory)
	err := svc.Register()
	require.NoError(t, err)
}

func TestComputeService_Start_Disabled(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, false)
	defer os.RemoveAll(tmpDir)

	err := svc.Start()
	require.NoError(t, err)
}

func TestComputeService_Stop_Disabled(t *testing.T) {
	svc, tmpDir := newTestComputeService(t, false)
	defer os.RemoveAll(tmpDir)

	err := svc.Stop()
	require.NoError(t, err)
}

func TestComputeService_RunnerDefinition_Sorting(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	runners := []types.RunnerDefinition{
		{Name: "third", Order: 30},
		{Name: "first", Order: 10},
		{Name: "second", Order: 20},
	}
	svc.setRunners(runners)

	sorted := svc.getRunners()
	assert.Equal(t, "first", sorted[0].Name)
	assert.Equal(t, "second", sorted[1].Name)
	assert.Equal(t, "third", sorted[2].Name)
}

func TestComputeService_ReverseRunners(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	runners := []types.RunnerDefinition{
		{Name: "first", Order: 10},
		{Name: "second", Order: 20},
		{Name: "third", Order: 30},
	}
	svc.setRunners(runners)

	reversed := svc.getReverseRunners()
	assert.Equal(t, "third", reversed[0].Name)
	assert.Equal(t, "second", reversed[1].Name)
	assert.Equal(t, "first", reversed[2].Name)
}

func TestComputeService_GetProvider_NotFound(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	svc.setProviders([]types.IBasePlateformProvider{})

	provider := svc.GetProvider("nonexistent")
	assert.Nil(t, provider)
}
