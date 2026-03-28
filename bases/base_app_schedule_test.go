package bases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tuxounet/k2-sdk/types"
)

func TestNewBaseAppSchedule(t *testing.T) {
	handler := func(fire string) string { return "done-" + fire }
	schedule := NewBaseAppSchedule("cleanup", "0 0 * * * *", handler)

	assert.Equal(t, "cleanup", schedule.GetName())
	assert.Equal(t, "0 0 * * * *", schedule.GetCron())
	assert.NotNil(t, schedule.GetTaskHandler())
}

func TestBaseAppSchedule_GetTaskHandler_Invocation(t *testing.T) {
	handler := func(fire string) string { return "result-" + fire }
	schedule := NewBaseAppSchedule("test", "* * * * *", handler)

	result := schedule.GetTaskHandler()("2024-01-01")
	assert.Equal(t, "result-2024-01-01", result)
}

func TestBaseAppSchedule_NilHandler(t *testing.T) {
	var handler types.AppScheduleHandler = nil
	schedule := NewBaseAppSchedule("empty", "* * * * *", handler)

	assert.Nil(t, schedule.GetTaskHandler())
}
