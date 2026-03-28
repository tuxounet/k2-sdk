package scheduler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/testutils"
	"github.com/tuxounet/k2-sdk/types"
)

func newTestSchedulerService(t *testing.T) *Service {
	t.Helper()
	kernel := testutils.NewMockKernel("test-app", "1.0.0")

	tmpDir, err := os.MkdirTemp("", "k2-sched-test-*")
	require.NoError(t, err)
	kernel.SetRunDir(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	app := testutils.NewMockApp("test-app", "1.0.0")
	kernel.SetApp(app)

	svc := NewService(kernel).(*Service)
	return svc
}

func TestSchedulerService_NewService(t *testing.T) {
	svc := newTestSchedulerService(t)

	assert.Equal(t, ServiceKey, svc.GetName())
	assert.NotNil(t, svc.getCronRunner())
}

func TestSchedulerService_Init(t *testing.T) {
	svc := newTestSchedulerService(t)
	err := svc.Init()
	assert.NoError(t, err)
}

func TestSchedulerService_Start_NoSchedules(t *testing.T) {
	svc := newTestSchedulerService(t)

	err := svc.Start()
	require.NoError(t, err)
}

func TestSchedulerService_Start_WithSchedules(t *testing.T) {
	kernel := testutils.NewMockKernel("test-app", "1.0.0")
	app := testutils.NewMockApp("test-app", "1.0.0")

	// Create a mock controller with a schedule
	comp := testutils.NewMockComponent(app, "comp", 0)
	ctrl := testutils.NewMockController(comp, "ctrl", 0)
	ctrl.AddTestSchedule("test-sched", "0 0 0 1 1 *", func(fire string) string { return "ok" })
	comp.SetControllers([]types.IAppController{ctrl})
	app.SetComponents([]types.IAppComponent{comp})
	kernel.SetApp(app)

	svc := NewService(kernel).(*Service)

	err := svc.Start()
	require.NoError(t, err)

	err = svc.Stop()
	require.NoError(t, err)
}

func TestSchedulerService_Stop(t *testing.T) {
	svc := newTestSchedulerService(t)

	err := svc.Start()
	require.NoError(t, err)

	err = svc.Stop()
	require.NoError(t, err)
}
