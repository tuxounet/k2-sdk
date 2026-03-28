package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

// --- MockStoreProvider ---

type mockStoreProvider struct {
	name string
	data map[string][]byte
}

func newMockStoreProvider(name string) *mockStoreProvider {
	return &mockStoreProvider{name: name, data: make(map[string][]byte)}
}

func (m *mockStoreProvider) GetName() string                 { return m.name }
func (m *mockStoreProvider) Setup() error                    { return nil }
func (m *mockStoreProvider) GetLogger() runtimeTypes.ILogger { return nil }
func (m *mockStoreProvider) Exists(store *Store, path string) (bool, error) {
	_, ok := m.data[path]
	return ok, nil
}
func (m *mockStoreProvider) Read(store *Store, path string) ([]byte, error) {
	return m.data[path], nil
}
func (m *mockStoreProvider) Write(store *Store, path string, data []byte) error {
	m.data[path] = data
	return nil
}
func (m *mockStoreProvider) Delete(store *Store, path string) error {
	delete(m.data, path)
	return nil
}

// --- Tests ---

func TestNewStore(t *testing.T) {
	flags := map[string]interface{}{"path": "/tmp/test"}
	store := NewStore("test-store", "local", flags)

	assert.Equal(t, "test-store", store.GetName())
	assert.Equal(t, "local", store.Backend)
	assert.Equal(t, "/tmp/test", store.Flags["path"])
}

func TestStore_ResolveBackend_Found(t *testing.T) {
	store := NewStore("test", "mock", nil)
	provider := newMockStoreProvider("mock")

	err := store.ResolveBackend([]IStoreProvider{provider})
	require.NoError(t, err)
}

func TestStore_ResolveBackend_NotFound(t *testing.T) {
	store := NewStore("test", "unknown", nil)
	provider := newMockStoreProvider("mock")

	err := store.ResolveBackend([]IStoreProvider{provider})
	// ResolveBackend returns nil even if no matching backend found
	assert.NoError(t, err)
}

func TestStore_Operations_WithBackend(t *testing.T) {
	store := NewStore("test", "mock", nil)
	provider := newMockStoreProvider("mock")

	err := store.ResolveBackend([]IStoreProvider{provider})
	require.NoError(t, err)

	// Write
	err = store.WriteObject("key1", []byte("hello"))
	require.NoError(t, err)

	// Exists
	exists, err := store.Exists("key1")
	require.NoError(t, err)
	assert.True(t, exists)

	// Read
	data, err := store.ReadObject("key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), data)

	// Delete
	err = store.DeleteObject("key1")
	require.NoError(t, err)

	exists, err = store.Exists("key1")
	require.NoError(t, err)
	assert.False(t, exists)
}
