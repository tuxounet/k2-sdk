package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunnerVerbConstants(t *testing.T) {
	assert.Equal(t, RunnerVerb("provision"), RunnerVerbProvision)
	assert.Equal(t, RunnerVerb("start"), RunnerVerbStart)
	assert.Equal(t, RunnerVerb("stop"), RunnerVerbStop)
	assert.Equal(t, RunnerVerb("teardown"), RunnerVerbTeardown)
}

func TestRunnerDefinition_Construction(t *testing.T) {
	rd := RunnerDefinition{
		Order:     1,
		Provider:  "container",
		Name:      "web-server",
		Provision: "provision-web",
		Teardown:  "teardown-web",
		Start:     "start-web",
		Stop:      "stop-web",
	}

	assert.Equal(t, 1, rd.Order)
	assert.Equal(t, "container", rd.Provider)
	assert.Equal(t, "web-server", rd.Name)
	assert.Equal(t, "provision-web", rd.Provision)
	assert.Equal(t, "teardown-web", rd.Teardown)
	assert.Equal(t, "start-web", rd.Start)
	assert.Equal(t, "stop-web", rd.Stop)
}

func TestRunnerDefinition_ZeroValue(t *testing.T) {
	rd := RunnerDefinition{}

	assert.Equal(t, 0, rd.Order)
	assert.Equal(t, "", rd.Provider)
	assert.Equal(t, "", rd.Name)
}

func TestRunnerVerb_StringConversion(t *testing.T) {
	assert.Equal(t, "provision", string(RunnerVerbProvision))
	assert.Equal(t, "start", string(RunnerVerbStart))
	assert.Equal(t, "stop", string(RunnerVerbStop))
	assert.Equal(t, "teardown", string(RunnerVerbTeardown))
}
