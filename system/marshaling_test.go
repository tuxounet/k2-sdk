package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- LoadYamlFromString ---

func TestLoadYamlFromString_ValidYAML(t *testing.T) {
	yaml := `name: test
version: "1.0"
count: 42`

	type Config struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Count   int    `yaml:"count"`
	}

	result, err := LoadYamlFromString[Config](yaml)
	require.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, "1.0", result.Version)
	assert.Equal(t, 42, result.Count)
}

func TestLoadYamlFromString_Map(t *testing.T) {
	yaml := `host:
  port: 8080
  name: localhost`

	result, err := LoadYamlFromString[map[string]any](yaml)
	require.NoError(t, err)

	host := result["host"].(map[string]any)
	assert.Equal(t, 8080, host["port"])
	assert.Equal(t, "localhost", host["name"])
}

func TestLoadYamlFromString_InvalidYAML(t *testing.T) {
	yaml := `invalid: [unclosed`

	_, err := LoadYamlFromString[map[string]any](yaml)
	assert.Error(t, err)
}

func TestLoadYamlFromString_EmptyString(t *testing.T) {
	result, err := LoadYamlFromString[map[string]any]("")
	require.NoError(t, err)
	assert.Nil(t, result)
}

// --- LoadJSONFromString ---

func TestLoadJSONFromString_ValidJSON(t *testing.T) {
	json := `{"name":"test","count":42}`

	type Config struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	result, err := LoadJSONFromString[Config](json)
	require.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Count)
}

func TestLoadJSONFromString_NestedJSON(t *testing.T) {
	json := `{"host":{"port":8080,"name":"localhost"}}`

	result, err := LoadJSONFromString[map[string]any](json)
	require.NoError(t, err)

	host := result["host"].(map[string]any)
	assert.Equal(t, float64(8080), host["port"])
	assert.Equal(t, "localhost", host["name"])
}

func TestLoadJSONFromString_InvalidJSON(t *testing.T) {
	json := `{invalid}`

	_, err := LoadJSONFromString[map[string]any](json)
	assert.Error(t, err)
}

// --- DumpToYamlString ---

func TestDumpToYamlString_Map(t *testing.T) {
	data := map[string]any{
		"name":  "test",
		"count": 42,
	}

	result, err := DumpToYamlString(data)
	require.NoError(t, err)
	assert.Contains(t, result, "name: test")
	assert.Contains(t, result, "count: 42")
}

func TestDumpToYamlString_Struct(t *testing.T) {
	type Config struct {
		Name string `yaml:"name"`
	}

	result, err := DumpToYamlString(Config{Name: "hello"})
	require.NoError(t, err)
	assert.Contains(t, result, "name: hello")
}

// --- DumpToJsonString ---

func TestDumpToJsonString_Map(t *testing.T) {
	data := map[string]any{
		"name":  "test",
		"count": 42,
	}

	result, err := DumpToJsonString(data)
	require.NoError(t, err)
	assert.Contains(t, result, `"name":"test"`)
	assert.Contains(t, result, `"count":42`)
}

func TestDumpToJsonString_Struct(t *testing.T) {
	type Config struct {
		Name string `json:"name"`
	}

	result, err := DumpToJsonString(Config{Name: "hello"})
	require.NoError(t, err)
	assert.Equal(t, `{"name":"hello"}`, result)
}

func TestDumpToJsonString_Nil(t *testing.T) {
	result, err := DumpToJsonString(nil)
	require.NoError(t, err)
	assert.Equal(t, "null", result)
}

// --- Roundtrip ---

func TestYamlRoundtrip(t *testing.T) {
	type Config struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	original := Config{Name: "test", Version: "1.0"}
	yamlStr, err := DumpToYamlString(original)
	require.NoError(t, err)

	restored, err := LoadYamlFromString[Config](yamlStr)
	require.NoError(t, err)
	assert.Equal(t, original, restored)
}

func TestJsonRoundtrip(t *testing.T) {
	type Config struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	original := Config{Name: "test", Version: "1.0"}
	jsonStr, err := DumpToJsonString(original)
	require.NoError(t, err)

	restored, err := LoadJSONFromString[Config](jsonStr)
	require.NoError(t, err)
	assert.Equal(t, original, restored)
}
