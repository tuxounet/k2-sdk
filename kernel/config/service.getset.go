package config

import (
	"fmt"
	"strings"
)

func (c *Service) GetCurrent() map[string]any {
	return c.GetData("records").(map[string]any)
}

func (c *Service) Has(key string) bool {
	return c.Get(key) != nil
}

// Get retrieves the value at the given path from the map.
func (c *Service) Get(path string) any {
	m := c.GetCurrent()
	keys := strings.Split(path, ".")
	var result any = m
	for _, key := range keys {
		key = strings.TrimSpace(key)
		key = strings.ToLower(key)
		if resultMap, ok := result.(map[string]any); ok {
			ret, found := resultMap[key]
			if !found {
				return nil
			}
			result = ret

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

func (c *Service) GetAsIntOrDefault(key string, def int) (int, error) {
	value := c.Get(key)
	if value == nil {
		return def, nil
	}

	return value.(int), nil
}

func (c *Service) SetValue(key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	m := c.GetCurrent()
	keys := strings.Split(key, ".")
	var current any = m

	for i, k := range keys {
		if i == len(keys)-1 {
			if resultMap, ok := current.(map[string]any); ok {
				resultMap[k] = value
			} else {
				return fmt.Errorf("cannot set value at %s, not a map", key)
			}
		} else {
			if resultMap, ok := current.(map[string]any); ok {
				if next, exists := resultMap[k]; exists {
					current = next
				} else {
					nextMap := make(map[string]any)
					resultMap[k] = nextMap
					current = nextMap
				}
			} else {
				return fmt.Errorf("cannot traverse %s, not a map", key)
			}
		}
	}
	c.GetLogger().DebugF("Set value at %s to %v", key, value)

	c.SetData("records", m)
	return nil
}
