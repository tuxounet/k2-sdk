package bases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuxounet/k2-sdk/testutils"
	"github.com/tuxounet/k2-sdk/types"
)

func newTestController() *BaseAppController {
	app := testutils.NewMockApp("test-app", "1.0")
	component := testutils.NewMockComponent(app, "test-comp", 0)
	ctrl := NewBaseAppController(component, "test-ctrl", 1, nil, types.AccessPolicyPublic)
	return &ctrl
}

func TestNewBaseAppController(t *testing.T) {
	ctrl := newTestController()

	assert.Equal(t, "test-ctrl", ctrl.GetName())
	assert.Equal(t, 1, ctrl.GetOrder())
	assert.NotNil(t, ctrl.GetComponent())
	assert.NotNil(t, ctrl.GetLogger())
}

func TestBaseAppController_AccessPolicy(t *testing.T) {
	app := testutils.NewMockApp("app", "1.0")
	component := testutils.NewMockComponent(app, "comp", 0)

	ctrl := NewBaseAppController(component, "ctrl", 0, nil, types.AccessPolicyAuthenticated)
	assert.Equal(t, types.AccessPolicyAuthenticated, ctrl.GetAccessPolicy())
}

func TestBaseAppController_Data_GetSet(t *testing.T) {
	ctrl := newTestController()

	ctrl.SetData("key1", "value1")
	ctrl.SetData("key2", 42)

	assert.Equal(t, "value1", ctrl.GetData("key1"))
	assert.Equal(t, 42, ctrl.GetData("key2"))
}

func TestBaseAppController_Data_MissingKey(t *testing.T) {
	ctrl := newTestController()
	assert.Nil(t, ctrl.GetData("nonexistent"))
}

func TestBaseAppController_Data_Overwrite(t *testing.T) {
	ctrl := newTestController()

	ctrl.SetData("key", "original")
	ctrl.SetData("key", "updated")

	assert.Equal(t, "updated", ctrl.GetData("key"))
}

func TestBaseAppController_AddSchedule(t *testing.T) {
	ctrl := newTestController()

	assert.Empty(t, ctrl.GetSchedules())

	ctrl.AddSchedule("task1", "0 * * * * *", func(fire string) string { return "ok" })
	assert.Len(t, ctrl.GetSchedules(), 1)
	assert.Equal(t, "task1", ctrl.GetSchedules()[0].GetName())
	assert.Equal(t, "0 * * * * *", ctrl.GetSchedules()[0].GetCron())

	ctrl.AddSchedule("task2", "0 0 * * * *", func(fire string) string { return "ok2" })
	assert.Len(t, ctrl.GetSchedules(), 2)
}

func TestBaseAppController_Lifecycle(t *testing.T) {
	ctrl := newTestController()

	assert.NoError(t, ctrl.Init())
	assert.NoError(t, ctrl.Register(nil))
	assert.NoError(t, ctrl.Start())
	assert.NoError(t, ctrl.Stop())
}

func TestBaseAppController_Config_Nil(t *testing.T) {
	ctrl := newTestController()
	assert.Nil(t, ctrl.GetConfig())
}

func TestBaseAppController_Component(t *testing.T) {
	ctrl := newTestController()
	assert.Equal(t, "test-comp", ctrl.GetComponent().GetName())
}
