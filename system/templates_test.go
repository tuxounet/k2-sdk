package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnTemplateWithGoTemplate_SimpleSubstitution(t *testing.T) {
	tpl := "Hello {{ .Name }}"
	data := map[string]any{"Name": "World"}

	result, err := UnTemplateWithGoTemplate(tpl, data)
	require.NoError(t, err)
	assert.Equal(t, "Hello World", result)
}

func TestUnTemplateWithGoTemplate_MultipleFields(t *testing.T) {
	tpl := "{{ .Name }} v{{ .Version }}"
	data := map[string]any{"Name": "k2-sdk", "Version": "1.0"}

	result, err := UnTemplateWithGoTemplate(tpl, data)
	require.NoError(t, err)
	assert.Equal(t, "k2-sdk v1.0", result)
}

func TestUnTemplateWithGoTemplate_SprigFunction(t *testing.T) {
	tpl := "{{ .Name | upper }}"
	data := map[string]any{"Name": "hello"}

	result, err := UnTemplateWithGoTemplate(tpl, data)
	require.NoError(t, err)
	assert.Equal(t, "HELLO", result)
}

func TestUnTemplateWithGoTemplate_SprigDefault(t *testing.T) {
	tpl := `{{ .Missing | default "fallback" }}`
	data := map[string]any{}

	result, err := UnTemplateWithGoTemplate(tpl, data)
	require.NoError(t, err)
	assert.Equal(t, "fallback", result)
}

func TestUnTemplateWithGoTemplate_InvalidTemplate(t *testing.T) {
	tpl := "{{ .Name"
	data := map[string]any{"Name": "test"}

	_, err := UnTemplateWithGoTemplate(tpl, data)
	assert.Error(t, err)
}

func TestUnTemplateWithGoTemplate_EmptyTemplate(t *testing.T) {
	result, err := UnTemplateWithGoTemplate("", nil)
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestUnTemplateWithGoTemplate_NoSubstitution(t *testing.T) {
	tpl := "plain text no templates"
	result, err := UnTemplateWithGoTemplate(tpl, nil)
	require.NoError(t, err)
	assert.Equal(t, "plain text no templates", result)
}

func TestUnTemplateWithGoTemplate_NestedData(t *testing.T) {
	tpl := "{{ .Config.Port }}"
	data := map[string]any{
		"Config": map[string]any{"Port": 8080},
	}

	result, err := UnTemplateWithGoTemplate(tpl, data)
	require.NoError(t, err)
	assert.Equal(t, "8080", result)
}
