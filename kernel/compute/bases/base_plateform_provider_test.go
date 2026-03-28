package bases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/kernel/compute/types"
	"github.com/tuxounet/k2-sdk/testutils"
)

func newTestProvider(t *testing.T) *BasePlateformProvider[string] {
	t.Helper()
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := testutils.NewMockKernelService(kernel, "compute")
	base := NewBasePlateformProvider[string](svc, "test-provider", nil)
	return &base
}

func TestBasePlateformProvider_NewBasePlateformProvider(t *testing.T) {
	p := newTestProvider(t)
	assert.Equal(t, "test-provider", p.GetName())
	assert.NotNil(t, p.GetLogger())
	assert.NotNil(t, p.GetService())
}

func TestBasePlateformProvider_IsRequired(t *testing.T) {
	p := newTestProvider(t)

	assert.False(t, p.GetIsRequired())
	p.SetIsRequired(true)
	assert.True(t, p.GetIsRequired())
}

func TestBasePlateformProvider_Data(t *testing.T) {
	p := newTestProvider(t)

	assert.Nil(t, p.GetData("key"))
	p.SetData("key", "value")
	assert.Equal(t, "value", p.GetData("key"))
}

func TestBasePlateformProvider_Definitions(t *testing.T) {
	p := newTestProvider(t)

	err := p.RegisterDefinition("def1")
	require.NoError(t, err)
	err = p.RegisterDefinition("def2")
	require.NoError(t, err)

	defs := p.GetDefinitions()
	assert.Len(t, defs, 2)
	assert.Equal(t, "def1", defs[0])
	assert.Equal(t, "def2", defs[1])
}

func TestBasePlateformProvider_ResetDefinitions(t *testing.T) {
	p := newTestProvider(t)

	p.RegisterDefinition("def1")
	p.ResetDefinitions()

	defs := p.GetDefinitions()
	assert.Empty(t, defs)
}

func TestBasePlateformProvider_Nuke(t *testing.T) {
	p := newTestProvider(t)

	p.RegisterDefinition("def1")
	p.RegisterDefinition("def2")

	err := p.Nuke()
	require.NoError(t, err)

	defs := p.GetDefinitions()
	assert.Empty(t, defs)
}

func TestBasePlateformProvider_Lifecycle(t *testing.T) {
	p := newTestProvider(t)

	assert.NoError(t, p.Init())
	assert.NoError(t, p.Setup())
	assert.NoError(t, p.Start())
	assert.NoError(t, p.Stop())
}

func TestBasePlateformProvider_Render(t *testing.T) {
	p := newTestProvider(t)

	runners, err := p.Render()
	assert.NoError(t, err)
	assert.Nil(t, runners)
}

func TestBasePlateformProvider_GetRunners(t *testing.T) {
	p := newTestProvider(t)

	runners := p.GetRunners()
	assert.Empty(t, runners)
}

func TestBasePlateformProvider_WithIngressRegistar(t *testing.T) {
	kernel := testutils.NewMockKernel("test", "1.0")
	svc := testutils.NewMockKernelService(kernel, "compute")

	called := false
	registar := func(ingress interface{}) error {
		called = true
		return nil
	}
	_ = registar
	_ = called

	base := NewBasePlateformProvider[types.RunnerDefinition](svc, "test-with-ingress", nil)
	assert.Equal(t, "test-with-ingress", base.GetName())
	assert.Nil(t, base.GetIngressRegistar())
}
