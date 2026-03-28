package bases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
	"github.com/tuxounet/k2-sdk/types"
)

func TestNewBaseApp(t *testing.T) {
	app := NewBaseApp("test-app", "1.0.0", nil, nil, nil, nil, nil)
	require.NotNil(t, app)

	assert.Equal(t, "test-app", app.GetName())
	assert.Equal(t, "1.0.0", app.GetVersion())
}

func TestNewBaseApp_TrimsWhitespace(t *testing.T) {
	app := NewBaseApp("  my-app  ", "  2.0  ", nil, nil, nil, nil, nil)

	assert.Equal(t, "my-app", app.GetName())
	assert.Equal(t, "2.0", app.GetVersion())
}

func TestNewBaseApp_NilResources(t *testing.T) {
	app := NewBaseApp("app", "1.0", nil, nil, nil, nil, nil)

	assert.Nil(t, app.GetDocs())
	assert.Nil(t, app.GetUI())
	assert.Nil(t, app.GetConfig())
	assert.Nil(t, app.GetExternals())
}

func TestBaseApp_KernelGetSet(t *testing.T) {
	app := NewBaseApp("app", "1.0", nil, nil, nil, nil, nil)

	assert.Nil(t, app.GetKernel())

	kernel := testutils.NewMockKernel("test", "1.0")
	app.SetKernel(kernel)
	assert.Equal(t, kernel, app.GetKernel())
}

func TestBaseApp_LoggerGetSet(t *testing.T) {
	app := NewBaseApp("app", "1.0", nil, nil, nil, nil, nil)

	assert.Nil(t, app.GetLogger())

	logger := testutils.NewMockLogger("test")
	app.SetLogger(logger)
	assert.Equal(t, logger, app.GetLogger())
}

func TestBaseApp_GetComponents_NilKernel(t *testing.T) {
	app := NewBaseApp("app", "1.0", nil, nil, nil, nil, nil)
	assert.Nil(t, app.GetComponents())
}

func TestBaseApp_GetComponents_LazyInit(t *testing.T) {
	ctorCalled := 0
	ctor := func(app types.IApp) types.IAppComponent {
		ctorCalled++
		return testutils.NewMockComponent(app, "comp-1", 0)
	}

	app := NewBaseApp("app", "1.0", nil, nil, nil, []types.AppComponentCtor{ctor}, nil)
	kernel := testutils.NewMockKernel("test", "1.0")
	app.SetKernel(kernel)

	assert.Equal(t, 0, ctorCalled)

	components := app.GetComponents()
	assert.Equal(t, 1, ctorCalled)
	assert.Len(t, components, 1)

	// Second call should not invoke constructors again
	components2 := app.GetComponents()
	assert.Equal(t, 1, ctorCalled)
	assert.Len(t, components2, 1)
}

func TestBaseApp_GetComponent_Found(t *testing.T) {
	ctor := func(app types.IApp) types.IAppComponent {
		return testutils.NewMockComponent(app, "my-comp", 0)
	}

	app := NewBaseApp("app", "1.0", nil, nil, nil, []types.AppComponentCtor{ctor}, nil)
	kernel := testutils.NewMockKernel("test", "1.0")
	app.SetKernel(kernel)

	comp := app.GetComponent("my-comp")
	require.NotNil(t, comp)
	assert.Equal(t, "my-comp", comp.GetName())
}

func TestBaseApp_GetComponent_NotFound(t *testing.T) {
	app := NewBaseApp("app", "1.0", nil, nil, nil, nil, nil)
	kernel := testutils.NewMockKernel("test", "1.0")
	app.SetKernel(kernel)

	comp := app.GetComponent("nonexistent")
	assert.Nil(t, comp)
}

func TestBaseApp_AddComponent(t *testing.T) {
	app := NewBaseApp("app", "1.0", nil, nil, nil, nil, nil)
	kernel := testutils.NewMockKernel("test", "1.0")
	app.SetKernel(kernel)

	ctor := func(a types.IApp) types.IAppComponent {
		return testutils.NewMockComponent(a, "added-comp", 0)
	}

	app.AddComponent(ctor)

	components := app.GetComponents()
	assert.Len(t, components, 1)
	assert.Equal(t, "added-comp", components[0].GetName())
}

func TestBaseApp_MultipleComponents(t *testing.T) {
	ctors := []types.AppComponentCtor{
		func(a types.IApp) types.IAppComponent { return testutils.NewMockComponent(a, "auth", 0) },
		func(a types.IApp) types.IAppComponent { return testutils.NewMockComponent(a, "api", 1) },
		func(a types.IApp) types.IAppComponent { return testutils.NewMockComponent(a, "admin", 2) },
	}

	app := NewBaseApp("app", "1.0", nil, nil, nil, ctors, nil)
	kernel := testutils.NewMockKernel("test", "1.0")
	app.SetKernel(kernel)

	components := app.GetComponents()
	assert.Len(t, components, 3)

	assert.NotNil(t, app.GetComponent("auth"))
	assert.NotNil(t, app.GetComponent("api"))
	assert.NotNil(t, app.GetComponent("admin"))
	assert.Nil(t, app.GetComponent("nonexistent"))
}
