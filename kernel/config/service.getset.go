package config

import (
	"fmt"
	"strings"
)

func (c *Service) GetCurrent() map[string]interface{} {
	return c.GetData("records").(map[string]interface{})
}

func (c *Service) Has(key string) bool {
	return c.Get(key) != nil
}

// Get retrieves the value at the given path from the map.
func (c *Service) Get(path string) interface{} {
	m := c.GetCurrent()
	keys := strings.Split(path, ".")
	var result interface{} = m
	for _, key := range keys {
		if resultMap, ok := result.(map[string]interface{}); ok {
			result = resultMap[key]
		} else {
			return nil
		}
	}
	return result
}

func (c *Service) GetAsString(key string) (string, error) {

	value := c.Get(key)
	if value == nil {
		return "", fmt.Errorf("key %s not found", key)
	}

	return value.(string), nil
}

func (c *Service) GetAsInt(key string) (int, error) {

	value := c.Get(key)
	if value == nil {
		return -1, fmt.Errorf("key %s not found", key)
	}

	return value.(int), nil
}

func (c *Service) GetAsBool(key string) (bool, error) {
	value := c.Get(key)
	if value == nil {
		return false, fmt.Errorf("key %s not found", key)
	}

	return value.(bool), nil
}

func (c *Service) GetAsStringOrDefault(key string, def string) (string, error) {
	value := c.Get(key)
	if value == nil {
		return def, nil
	}

	return value.(string), nil
}
