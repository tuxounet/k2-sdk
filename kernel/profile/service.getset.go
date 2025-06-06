package profile

import (
	"os"
	"path"

	"github.com/tuxounet/k2-sdk/kernel/profile/bases"
)

func (c *ProfileService) GetProfile() (*bases.Profile, error) {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			c.GetLogger().Warn("Profile not found")
			return nil, nil
		}
		c.GetLogger().ErrorF("Failed to read profile: %s", err.Error())
		return nil, err
	}

	return profile, nil
}

func (c *ProfileService) GetUserDirectory() (string, error) {
	runDir := c.GetKernel().GetRunDirectory()
	dataDir := path.Join(runDir, "home")

	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		c.GetLogger().ErrorF("Failed to create user directory: %s", err.Error())
		return "", err
	}

	return dataDir, nil
}

func (c *ProfileService) GetPublicProfile() (*bases.ProfilePublic, error) {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			c.GetLogger().Warn("Profile not found")
			return nil, nil
		}
		c.GetLogger().ErrorF("Failed to read profile: %s", err.Error())
		return nil, err
	}
	public := profile.Public()

	return public, nil
}

func (c *ProfileService) HasProperty(key string) (bool, error) {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			c.GetLogger().Warn("Profile not found")
			return false, nil
		}
		c.GetLogger().ErrorF("Failed to read profile: %s", err.Error())
		return false, err
	}

	_, ok := profile.Properties[key]
	return ok, nil
}

func (c *ProfileService) GetProperty(key string) (string, error) {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			c.GetLogger().Warn("Profile not found")
			return "", nil
		}
		c.GetLogger().ErrorF("Failed to read profile: %s", err.Error())
		return "", err
	}

	value, ok := profile.Properties[key]
	if !ok {
		return "", nil
	}

	return value, nil
}

func (c *ProfileService) HasSecret(key string) (bool, error) {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	_, ok := profile.Secrets[key]
	return ok, nil
}

func (c *ProfileService) GetSecret(key string) (string, error) {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	value, ok := profile.Secrets[key]
	if !ok {
		return "", nil
	}

	return value, nil
}

func (c *ProfileService) SetProperty(key string, value string) error {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			profile = &bases.Profile{
				Properties: make(map[string]string),
			}
		} else {
			c.GetLogger().ErrorF("Failed to read profile: %s", err.Error())
			return err
		}
	}

	profile.Properties[key] = value

	err = c.writeProfile(profile)
	if err != nil {
		c.GetLogger().ErrorF("Failed to write profile: %s", err.Error())
		return err
	}

	return nil
}

func (c *ProfileService) SetSecret(key string, value string) error {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			profile = &bases.Profile{
				Secrets: make(map[string]string),
			}
		} else {
			c.GetLogger().ErrorF("Failed to read profile: %s", err.Error())
			return err
		}
	}

	profile.Secrets[key] = value

	err = c.writeProfile(profile)
	if err != nil {
		c.GetLogger().ErrorF("Failed to write profile: %s", err.Error())
		return err
	}

	return nil
}
