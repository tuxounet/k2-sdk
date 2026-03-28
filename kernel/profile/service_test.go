package profile

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newTestProfileService(t *testing.T) (*ProfileService, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "k2-profile-test-*")
	require.NoError(t, err)

	kernel := testutils.NewMockKernel("test-app", "1.0.0")
	kernel.SetRunDir(tmpDir)

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	svc := NewService(kernel).(*ProfileService)
	return svc, tmpDir
}

func TestProfileService_Init_CreatesProfile(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	profile, err := svc.GetProfile()
	require.NoError(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, "test-app", profile.Name)
	assert.Equal(t, "1.0.0", profile.Version)
}

func TestProfileService_GetProfile_NonExistent(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	profile, err := svc.GetProfile()
	assert.Nil(t, err)
	assert.Nil(t, profile)
}

func TestProfileService_SetProperty_GetProperty(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	err = svc.SetProperty("greeting", "hello")
	require.NoError(t, err)

	val, err := svc.GetProperty("greeting")
	require.NoError(t, err)
	assert.Equal(t, "hello", val)
}

func TestProfileService_HasProperty(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	has, err := svc.HasProperty("missing")
	require.NoError(t, err)
	assert.False(t, has)

	svc.SetProperty("exists", "yes")

	has, err = svc.HasProperty("exists")
	require.NoError(t, err)
	assert.True(t, has)
}

func TestProfileService_SetSecret_GetSecret(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	err = svc.SetSecret("api_key", "secret123")
	require.NoError(t, err)

	val, err := svc.GetSecret("api_key")
	require.NoError(t, err)
	assert.Equal(t, "secret123", val)
}

func TestProfileService_HasSecret(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	has, err := svc.HasSecret("missing")
	require.NoError(t, err)
	assert.False(t, has)

	svc.SetSecret("exists", "yes")

	has, err = svc.HasSecret("exists")
	require.NoError(t, err)
	assert.True(t, has)
}

func TestProfileService_GetPublicProfile(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	svc.SetSecret("hidden", "value")

	public, err := svc.GetPublicProfile()
	require.NoError(t, err)
	require.NotNil(t, public)
	assert.Equal(t, "test-app", public.Name)
	assert.Equal(t, "1.0.0", public.Version)
}

func TestProfileService_GetUserDirectory(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	dir, err := svc.GetUserDirectory()
	require.NoError(t, err)
	assert.Contains(t, dir, "home")

	// Directory should exist
	info, err := os.Stat(dir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestProfileService_GetProperty_MissingKey(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	val, err := svc.GetProperty("nonexistent")
	require.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestProfileService_SetProperty_CreatesProfileIfMissing(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	// Don't call Init() - profile doesn't exist yet
	err := svc.SetProperty("key", "value")
	require.NoError(t, err)

	val, err := svc.GetProperty("key")
	require.NoError(t, err)
	assert.Equal(t, "value", val)
}

func TestProfileService_Init_UpdatesExistingProfile(t *testing.T) {
	svc, tmpDir := newTestProfileService(t)
	defer os.RemoveAll(tmpDir)

	// First init
	err := svc.Init()
	require.NoError(t, err)

	// Set some properties
	svc.SetProperty("custom", "data")

	// Second init should update name/version but keep properties
	err = svc.Init()
	require.NoError(t, err)

	val, err := svc.GetProperty("custom")
	require.NoError(t, err)
	assert.Equal(t, "data", val)
}
