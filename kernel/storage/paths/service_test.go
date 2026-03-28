package paths

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newTestService() *Service {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := NewService(kernel)
	return svc.(*Service)
}

func TestNewService(t *testing.T) {
	svc := newTestService()
	assert.Equal(t, ServiceKey, svc.GetName())
}

func TestCominePath_MultipleParts(t *testing.T) {
	svc := newTestService()
	result := svc.CominePath("home", "user", "docs")
	assert.Equal(t, "home/user/docs", result)
}

func TestCominePath_SinglePart(t *testing.T) {
	svc := newTestService()
	result := svc.CominePath("home")
	assert.Equal(t, "home", result)
}

func TestCominePath_WithSlashes(t *testing.T) {
	svc := newTestService()
	result := svc.CominePath("/home/", "/user/")
	assert.Equal(t, "/home/user", result)
}

func TestCominePath_Empty(t *testing.T) {
	svc := newTestService()
	result := svc.CominePath()
	assert.Equal(t, "", result)
}

func TestSplitPath_Simple(t *testing.T) {
	svc := newTestService()
	result := svc.SplitPath("home/user/docs")
	assert.Equal(t, []string{"home", "user", "docs"}, result)
}

func TestSplitPath_MultipleParts(t *testing.T) {
	svc := newTestService()
	result := svc.SplitPath("home", "user", "docs")
	assert.Equal(t, []string{"home", "user", "docs"}, result)
}

func TestSplitPath_AbsolutePath(t *testing.T) {
	svc := newTestService()
	result := svc.SplitPath("/home/user")
	assert.Equal(t, []string{"", "home", "user"}, result)
}

func TestGetBaseName(t *testing.T) {
	svc := newTestService()
	assert.Equal(t, "file.txt", svc.GetBaseName("/home/user/file.txt"))
	assert.Equal(t, "dir", svc.GetBaseName("/home/user/dir"))
	assert.Equal(t, "/", svc.GetBaseName("/"))
}

func TestGetDirName(t *testing.T) {
	svc := newTestService()
	assert.Equal(t, "/home/user", svc.GetDirName("/home/user/file.txt"))
	assert.Equal(t, ".", svc.GetDirName("file.txt"))
}

func TestGetExtName(t *testing.T) {
	svc := newTestService()
	assert.Equal(t, ".txt", svc.GetExtName("file.txt"))
	assert.Equal(t, ".yaml", svc.GetExtName("config.yaml"))
	assert.Equal(t, "", svc.GetExtName("Makefile"))
	assert.Equal(t, ".gz", svc.GetExtName("archive.tar.gz"))
}
