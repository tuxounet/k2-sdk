package bases

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKeyPair_Success(t *testing.T) {
	kp, err := NewKeyPair()
	require.NoError(t, err)
	assert.NotNil(t, kp)
	assert.NotEmpty(t, kp.PublicKey)
	assert.NotEmpty(t, kp.PrivateKey)
}

func TestNewKeyPair_ValidPrivateKey(t *testing.T) {
	kp, err := NewKeyPair()
	require.NoError(t, err)

	privKey, err := x509.ParseECPrivateKey(kp.PrivateKey)
	require.NoError(t, err)
	assert.NotNil(t, privKey)
	assert.Equal(t, "P-256", privKey.Curve.Params().Name)
}

func TestNewKeyPair_ValidPublicKey(t *testing.T) {
	kp, err := NewKeyPair()
	require.NoError(t, err)

	pubKey, err := x509.ParsePKIXPublicKey(kp.PublicKey)
	require.NoError(t, err)
	assert.NotNil(t, pubKey)
}

func TestNewKeyPair_UniqueKeys(t *testing.T) {
	kp1, err := NewKeyPair()
	require.NoError(t, err)

	kp2, err := NewKeyPair()
	require.NoError(t, err)

	assert.NotEqual(t, kp1.PrivateKey, kp2.PrivateKey)
	assert.NotEqual(t, kp1.PublicKey, kp2.PublicKey)
}

func TestNewKeyPair_KeyPairCorresponds(t *testing.T) {
	kp, err := NewKeyPair()
	require.NoError(t, err)

	privKey, err := x509.ParseECPrivateKey(kp.PrivateKey)
	require.NoError(t, err)

	// Marshal the public key derived from the private key and compare
	pubFromPriv, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	require.NoError(t, err)
	assert.Equal(t, kp.PublicKey, pubFromPriv)
}
