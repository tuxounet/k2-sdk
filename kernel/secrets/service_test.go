package secrets

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/profile"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newTestSecretsService(t *testing.T) (*SecretsService, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "k2-secrets-test-*")
	require.NoError(t, err)

	kernel := testutils.NewMockKernel("test-app", "1.0.0")
	kernel.SetRunDir(tmpDir)

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	// Create and init profile service (dependency)
	profileSvc := profile.NewService(kernel).(*profile.ProfileService)
	kernel.SetService(profileSvc)
	err = profileSvc.Init()
	require.NoError(t, err)

	secretsSvc := NewService(kernel).(*SecretsService)
	return secretsSvc, tmpDir
}

func TestSecretsService_Init_GeneratesSSHKeys(t *testing.T) {
	svc, tmpDir := newTestSecretsService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	// SSH keys should now exist
	privKey, err := svc.GetSecret("ssh.private")
	require.NoError(t, err)
	assert.NotEmpty(t, privKey)

	pubKey, err := svc.GetSecret("ssh.public")
	require.NoError(t, err)
	assert.NotEmpty(t, pubKey)
}

func TestSecretsService_Init_DoesNotRegenerateKeys(t *testing.T) {
	svc, tmpDir := newTestSecretsService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	privKey1, _ := svc.GetSecret("ssh.private")

	// Second init should not regenerate
	err = svc.Init()
	require.NoError(t, err)

	privKey2, _ := svc.GetSecret("ssh.private")
	assert.Equal(t, privKey1, privKey2)
}

func TestSecretsService_SetSecret_GetSecret(t *testing.T) {
	svc, tmpDir := newTestSecretsService(t)
	defer os.RemoveAll(tmpDir)

	data := []byte("my-secret-data")
	err := svc.SetSecret("test-key", data)
	require.NoError(t, err)

	retrieved, err := svc.GetSecret("test-key")
	require.NoError(t, err)
	assert.Equal(t, data, retrieved)
}

func TestSecretsService_GetSecret_NotFound(t *testing.T) {
	svc, tmpDir := newTestSecretsService(t)
	defer os.RemoveAll(tmpDir)

	_, err := svc.GetSecret("nonexistent")
	assert.Error(t, err)
}

func TestSecretsService_SetSecret_Base64Encoding(t *testing.T) {
	svc, tmpDir := newTestSecretsService(t)
	defer os.RemoveAll(tmpDir)

	data := []byte("hello world")
	err := svc.SetSecret("encoded", data)
	require.NoError(t, err)

	// Verify the profile stores base64-encoded value
	profileSvc := svc.getProfileService()
	raw, err := profileSvc.GetSecret("encoded")
	require.NoError(t, err)

	decoded, err := base64.StdEncoding.DecodeString(raw)
	require.NoError(t, err)
	assert.Equal(t, data, decoded)
}
