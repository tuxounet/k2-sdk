package secrets

import (
	"encoding/base64"
	"errors"
)

func (c *SecretsService) GetSecret(key string) ([]byte, error) {
	profiles := c.getProfileService()
	hasProperty, err := profiles.HasSecret(key)
	if err != nil {
		c.GetLogger().ErrorF("Error checking secret %s: %s", key, err.Error())
		return nil, err
	}

	if !hasProperty {
		return nil, errors.New("Secret not found : " + key)
	}

	rawValue, err := profiles.GetSecret(key)
	if err != nil {
		c.GetLogger().ErrorF("Error reading secret in profile %s: %s", key, err.Error())
		return nil, err
	}

	value, err := base64.StdEncoding.DecodeString(rawValue)
	if err != nil {
		c.GetLogger().ErrorF("Error decoding secret %s: %s", key, err.Error())
		return nil, err
	}

	return value, nil
}

func (c *SecretsService) SetSecret(key string, value []byte) error {

	profiles := c.getProfileService()

	rawValue := base64.StdEncoding.EncodeToString(value)

	err := profiles.SetSecret(key, rawValue)
	if err != nil {
		c.GetLogger().ErrorF("Error setting secret %s: %s", key, err.Error())
		return err
	}

	return nil
}
