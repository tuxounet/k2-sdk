package stores

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/profile"
	storesBases "github.com/tuxounet/k2-sdk/kernel/storage/stores/bases"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newTestStoreService(t *testing.T) (*Service, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "k2-stores-test-*")
	require.NoError(t, err)

	kernel := testutils.NewMockKernel("test-app", "1.0.0")
	kernel.SetRunDir(tmpDir)

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	// Create and init profile service (dependency for stores)
	profileSvc := profile.NewService(kernel).(*profile.ProfileService)
	kernel.SetService(profileSvc)
	err = profileSvc.Init()
	require.NoError(t, err)

	svc := NewService(kernel).(*Service)
	return svc, tmpDir
}

func TestStoreService_NewService(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)

	assert.Equal(t, ServiceKey, svc.GetName())
	assert.Equal(t, kernel, svc.GetKernel())
}

func TestStoreService_Init(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	// After init, root and local stores should exist
	stores := svc.getStores()
	assert.GreaterOrEqual(t, len(stores), 2)

	names, err := svc.GetStores()
	require.NoError(t, err)
	assert.Contains(t, names, "root")
	assert.Contains(t, names, "local")
}

func TestStoreService_Init_CreatesBackends(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	backends, err := svc.GetBackends()
	require.NoError(t, err)
	require.NotNil(t, backends)
	assert.Len(t, backends, 2) // local + rclone

	backendNames := make([]string, len(backends))
	for i, b := range backends {
		backendNames[i] = b.GetName()
	}
	assert.Contains(t, backendNames, "local")
	assert.Contains(t, backendNames, "rclone")
}

func TestStoreService_GetStore_Root(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	store, err := svc.GetStore("root")
	require.NoError(t, err)
	require.NotNil(t, store)
	assert.Equal(t, "root", store.GetName())
}

func TestStoreService_GetStore_Local(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	store, err := svc.GetStore("local")
	require.NoError(t, err)
	require.NotNil(t, store)
	assert.Equal(t, "local", store.GetName())
}

func TestStoreService_GetStore_NotFound(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	store, err := svc.GetStore("nonexistent")
	require.NoError(t, err)
	assert.Nil(t, store)
}

func TestStoreService_UpsertStore_CreateNew(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	newStore := &types.Store{
		Name:    "custom",
		Backend: "local",
		Flags:   map[string]any{"path": tmpDir},
	}
	err = svc.UpsertStore(newStore)
	require.NoError(t, err)

	names, err := svc.GetStores()
	require.NoError(t, err)
	assert.Contains(t, names, "custom")
}

func TestStoreService_UpsertStore_UpdateExisting(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	// Update root store with new flags
	updatedStore := &types.Store{
		Name:    "root",
		Backend: "local",
		Flags:   map[string]any{"path": tmpDir, "extra": "value"},
	}
	err = svc.UpsertStore(updatedStore)
	require.NoError(t, err)

	store, err := svc.GetStore("root")
	require.NoError(t, err)
	require.NotNil(t, store)
	assert.Equal(t, "value", store.Flags["extra"])
}

func TestStoreService_RootStore_WriteReadDelete(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	store, err := svc.GetStore("root")
	require.NoError(t, err)

	// Write
	err = store.WriteObject("test/data.txt", []byte("hello world"))
	require.NoError(t, err)

	// Exists
	exists, err := store.Exists("test/data.txt")
	require.NoError(t, err)
	assert.True(t, exists)

	// Read
	data, err := store.ReadObject("test/data.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("hello world"), data)

	// Delete
	err = store.DeleteObject("test/data.txt")
	require.NoError(t, err)

	exists, err = store.Exists("test/data.txt")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestStoreService_Init_CleansTmpDir(t *testing.T) {
	svc, tmpDir := newTestStoreService(t)
	defer os.RemoveAll(tmpDir)

	err := svc.Init()
	require.NoError(t, err)

	// tmp/.keep should exist
	store, err := svc.GetStore("root")
	require.NoError(t, err)

	exists, err := store.Exists("tmp/.keep")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestStoreService_GetStores_Empty(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	names, err := svc.GetStores()
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestStoreService_GetBackends_NoInit(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel).(*Service)

	backends, err := svc.GetBackends()
	require.NoError(t, err)
	assert.Nil(t, backends)
}

// --- BaseObjectStore Integration Tests ---

func setupObjectStore(t *testing.T) (types.IBaseObjectStore[map[string]string], string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "k2-objstore-test-*")
	require.NoError(t, err)

	kernel := testutils.NewMockKernel("test-app", "1.0.0")
	kernel.SetRunDir(tmpDir)

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	profileSvc := profile.NewService(kernel).(*profile.ProfileService)
	kernel.SetService(profileSvc)
	err = profileSvc.Init()
	require.NoError(t, err)

	storesSvc := NewService(kernel).(*Service)
	kernel.SetService(storesSvc)
	err = storesSvc.Init()
	require.NoError(t, err)

	defaultValue := `{}`
	mockSvc := testutils.NewMockKernelService(kernel, "test")
	store := storesBases.NewObjectStore[map[string]string](kernel, mockSvc, "root", "test-data.json", defaultValue)

	return store, tmpDir
}

func TestBaseObjectStore_HasValue_Initially(t *testing.T) {
	store, tmpDir := setupObjectStore(t)
	defer os.RemoveAll(tmpDir)

	has, err := store.HasValue()
	require.NoError(t, err)
	assert.False(t, has)
}

func TestBaseObjectStore_SetValue_GetValue(t *testing.T) {
	store, tmpDir := setupObjectStore(t)
	defer os.RemoveAll(tmpDir)

	data := map[string]string{"key1": "value1", "key2": "value2"}
	err := store.SetValue(data)
	require.NoError(t, err)

	has, err := store.HasValue()
	require.NoError(t, err)
	assert.True(t, has)

	result, err := store.GetValue()
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "value1", (*result)["key1"])
	assert.Equal(t, "value2", (*result)["key2"])
}

func TestBaseObjectStore_GetValue_Default(t *testing.T) {
	store, tmpDir := setupObjectStore(t)
	defer os.RemoveAll(tmpDir)

	result, err := store.GetValue()
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, *result)
}

func TestBaseObjectStore_Nuke(t *testing.T) {
	store, tmpDir := setupObjectStore(t)
	defer os.RemoveAll(tmpDir)

	data := map[string]string{"key": "val"}
	err := store.SetValue(data)
	require.NoError(t, err)

	err = store.Nuke()
	require.NoError(t, err)

	has, err := store.HasValue()
	require.NoError(t, err)
	assert.False(t, has)
}

func TestBaseObjectStore_Nuke_NonExistent(t *testing.T) {
	store, tmpDir := setupObjectStore(t)
	defer os.RemoveAll(tmpDir)

	err := store.Nuke()
	require.NoError(t, err)
}
