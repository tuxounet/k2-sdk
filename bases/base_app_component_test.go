package bases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
	"github.com/tuxounet/k2-sdk/types"
)

func TestNewBaseAppComponent(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")
	component := NewBaseAppComponent(app, "my-comp", 5, nil, nil, nil, types.AccessPolicyPublic, nil)

	assert.Equal(t, "my-comp", component.GetName())
	assert.Equal(t, 5, component.GetOrder())
	assert.Equal(t, app, component.GetApp())
	assert.Equal(t, types.AccessPolicyPublic, component.GetAccessPolicy())
}

func TestBaseAppComponent_LazyControllerInit(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	ctorCalled := 0
	ctor := func(comp types.IAppComponent) types.IAppController {
		ctorCalled++
		return &BaseAppController{
			name:      "ctrl-1",
			order:     0,
			component: comp,
			log:       comp.GetLogger().CreateSubLogger("ctrl-1"),
			schedules: make([]types.IAppSchedule, 0),
			data:      make(map[string]any),
		}
	}

	component := NewBaseAppComponent(app, "comp", 0, nil, nil, nil, types.AccessPolicyPublic,
		[]types.AppControllerCtor{ctor})

	assert.Equal(t, 0, ctorCalled)

	controllers := component.GetControllers()
	assert.Equal(t, 1, ctorCalled)
	assert.Len(t, controllers, 1)

	// Second call should not invoke constructors again
	controllers2 := component.GetControllers()
	assert.Equal(t, 1, ctorCalled)
	assert.Len(t, controllers2, 1)
}

func TestBaseAppComponent_GetController_Found(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	ctor := func(comp types.IAppComponent) types.IAppController {
		return &BaseAppController{
			name:      "found-ctrl",
			component: comp,
			log:       comp.GetLogger().CreateSubLogger("found-ctrl"),
			schedules: make([]types.IAppSchedule, 0),
			data:      make(map[string]any),
		}
	}

	component := NewBaseAppComponent(app, "comp", 0, nil, nil, nil, types.AccessPolicyPublic,
		[]types.AppControllerCtor{ctor})

	ctrl := component.GetController("found-ctrl")
	require.NotNil(t, ctrl)
	assert.Equal(t, "found-ctrl", ctrl.GetName())
}

func TestBaseAppComponent_GetController_NotFound(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	component := NewBaseAppComponent(app, "comp", 0, nil, nil, nil, types.AccessPolicyPublic, nil)
	ctrl := component.GetController("nonexistent")
	assert.Nil(t, ctrl)
}

func TestBaseAppComponent_Lifecycle(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")
	component := NewBaseAppComponent(app, "comp", 0, nil, nil, nil, types.AccessPolicyPublic, nil)

	assert.NoError(t, component.Init())
	assert.NoError(t, component.Register(nil))
	assert.NoError(t, component.Start())
	assert.NoError(t, component.Stop())
}

func TestBaseAppComponent_NilResources(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")
	component := NewBaseAppComponent(app, "comp", 0, nil, nil, nil, types.AccessPolicyPublic, nil)

	assert.Nil(t, component.GetDocs())
	assert.Nil(t, component.GetUI())
	assert.Nil(t, component.GetConfig())
}

func TestBaseAppComponent_Logger(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")
	component := NewBaseAppComponent(app, "my-comp", 0, nil, nil, nil, types.AccessPolicyPublic, nil)

	assert.NotNil(t, component.GetLogger())
}

func TestBaseAppComponent_MultipleControllers(t *testing.T) {
	app := testutils.NewMockApp("test-app", "1.0")

	makeCtor := func(name string) types.AppControllerCtor {
		return func(comp types.IAppComponent) types.IAppController {
			return &BaseAppController{
				name:      name,
				component: comp,
				log:       comp.GetLogger().CreateSubLogger(name),
				schedules: make([]types.IAppSchedule, 0),
				data:      make(map[string]any),
			}
		}
	}

	component := NewBaseAppComponent(app, "comp", 0, nil, nil, nil, types.AccessPolicyPublic,
		[]types.AppControllerCtor{makeCtor("ctrl-a"), makeCtor("ctrl-b"), makeCtor("ctrl-c")})

	controllers := component.GetControllers()
	assert.Len(t, controllers, 3)

	assert.NotNil(t, component.GetController("ctrl-a"))
	assert.NotNil(t, component.GetController("ctrl-b"))
	assert.NotNil(t, component.GetController("ctrl-c"))
	assert.Nil(t, component.GetController("ctrl-d"))
}
