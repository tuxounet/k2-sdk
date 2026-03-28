package bases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfile_Public(t *testing.T) {
	profile := &Profile{
		Name:       "my-app",
		Version:    "1.0.0",
		Properties: map[string]string{"key": "value"},
		Secrets:    map[string]string{"secret": "hidden"},
	}

	public := profile.Public()
	assert.Equal(t, "my-app", public.Name)
	assert.Equal(t, "1.0.0", public.Version)
}

func TestProfile_ZeroValue(t *testing.T) {
	profile := &Profile{}
	assert.Equal(t, "", profile.Name)
	assert.Equal(t, "", profile.Version)
	assert.Nil(t, profile.Properties)
	assert.Nil(t, profile.Secrets)
}

func TestProfilePublic_ZeroValue(t *testing.T) {
	public := &ProfilePublic{}
	assert.Equal(t, "", public.Name)
	assert.Equal(t, "", public.Version)
}
