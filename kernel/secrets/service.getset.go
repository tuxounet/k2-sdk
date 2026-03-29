package secrets

import (
	"encoding/base64"
	"errors"
	"fmt"
)

func (c *SecretsService) GetSecret(key string) ([]byte, error) {
	profiles := c.getProfileService()
	hasProperty, err := profiles.HasSecret(key)
	if err != nil {
		c.GetLogger().ErrorF("error checking secret %s: %s", key, err.Error())
		return nil, fmt.Errorf("secrets: error checking existence of %s: %w", key, err)
	}

	if !hasProperty {
		return nil, errors.New("Secret not found : " + key)
	}

	rawValue, err := profiles.GetSecret(key)
	if err != nil {
		c.GetLogger().ErrorF("error reading secret %s from profile: %s", key, err.Error())
		return nil, fmt.Errorf("secrets: error reading %s from profile: %w", key, err)
	}

	value, err := base64.StdEncoding.DecodeString(rawValue)
	if err != nil {
		c.GetLogger().ErrorF("error decoding secret %s (base64): %s", key, err.Error())
		return nil, fmt.Errorf("secrets: error decoding %s (base64): %w", key, err)
	}

	return value, nil
}

func (c *SecretsService) SetSecret(key string, value []byte) error {

	profiles := c.getProfileService()

	rawValue := base64.StdEncoding.EncodeToString(value)

	err := profiles.SetSecret(key, rawValue)
	if err != nil {
		c.GetLogger().ErrorF("error setting secret %s: %s", key, err.Error())
		return fmt.Errorf("secrets: error setting %s in profile: %w", key, err)
	}

	return nil
}
