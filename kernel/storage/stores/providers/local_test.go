package providers

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newLocalProvider(t *testing.T) *Local {
	t.Helper()
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := testutils.NewMockKernelService(kernel, "storage.stores")
	return NewLocal(svc)
}

func newStoreWithPath(tmpDir string) *types.Store {
	return &types.Store{
		Name:    "test",
		Backend: "local",
		Flags:   map[string]any{"path": tmpDir},
	}
}

func TestLocal_NewLocal(t *testing.T) {
	local := newLocalProvider(t)
	assert.Equal(t, "local", local.GetName())
}

func TestLocal_Setup(t *testing.T) {
	local := newLocalProvider(t)
	err := local.Setup()
	assert.NoError(t, err)
}

func TestLocal_Write_Read(t *testing.T) {
	local := newLocalProvider(t)

	tmpDir, err := os.MkdirTemp("", "k2-local-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := newStoreWithPath(tmpDir)

	err = local.Write(store, "test.txt", []byte("hello"))
	require.NoError(t, err)

	data, err := local.Read(store, "test.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), data)
}

func TestLocal_Exists(t *testing.T) {
	local := newLocalProvider(t)

	tmpDir, err := os.MkdirTemp("", "k2-local-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := newStoreWithPath(tmpDir)

	exists, err := local.Exists(store, "nonexistent.txt")
	require.NoError(t, err)
	assert.False(t, exists)

	err = local.Write(store, "exists.txt", []byte("data"))
	require.NoError(t, err)

	exists, err = local.Exists(store, "exists.txt")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestLocal_Delete_File(t *testing.T) {
	local := newLocalProvider(t)

	tmpDir, err := os.MkdirTemp("", "k2-local-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := newStoreWithPath(tmpDir)

	err = local.Write(store, "to-delete.txt", []byte("bye"))
	require.NoError(t, err)

	err = local.Delete(store, "to-delete.txt")
	require.NoError(t, err)

	exists, err := local.Exists(store, "to-delete.txt")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestLocal_Delete_Directory(t *testing.T) {
	local := newLocalProvider(t)

	tmpDir, err := os.MkdirTemp("", "k2-local-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := newStoreWithPath(tmpDir)

	err = local.Write(store, "subdir/file.txt", []byte("data"))
	require.NoError(t, err)

	err = local.Delete(store, "subdir")
	require.NoError(t, err)

	exists, err := local.Exists(store, "subdir")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestLocal_Delete_NonExistent(t *testing.T) {
	local := newLocalProvider(t)

	tmpDir, err := os.MkdirTemp("", "k2-local-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := newStoreWithPath(tmpDir)

	err = local.Delete(store, "nonexistent.txt")
	require.NoError(t, err)
}

func TestLocal_Write_CreatesDirectories(t *testing.T) {
	local := newLocalProvider(t)

	tmpDir, err := os.MkdirTemp("", "k2-local-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store := newStoreWithPath(tmpDir)

	err = local.Write(store, "a/b/c/deep.txt", []byte("deep"))
	require.NoError(t, err)

	content, err := os.ReadFile(path.Join(tmpDir, "a", "b", "c", "deep.txt"))
	require.NoError(t, err)
	assert.Equal(t, []byte("deep"), content)
}

func TestRClone_NewRClone(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := testutils.NewMockKernelService(kernel, "storage.stores")
	rclone := NewRClone(svc)
	assert.Equal(t, "rclone", rclone.GetName())
}

func TestRClone_Setup(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := testutils.NewMockKernelService(kernel, "storage.stores")
	rclone := NewRClone(svc)
	err := rclone.Setup()
	assert.NoError(t, err)
}
