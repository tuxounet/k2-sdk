package profile

import (
	"os"
	"path"

	"github.com/tuxounet/k2-sdk/kernel/profile/bases"
	"github.com/tuxounet/k2-sdk/system"
)

func (c *ProfileService) getProfileFilename() (string, error) {
	runDir := c.GetKernel().GetRunDirectory()
	etcDir := path.Join(runDir, "etc")

	err := os.MkdirAll(etcDir, 0755)
	if err != nil {
		c.GetLogger().ErrorF("Failed to create etc directory: %s", err.Error())
		return "", err
	}
	profileFile := path.Join(etcDir, "profile.json")

	return profileFile, nil
}

func (c *ProfileService) readProfile() (*bases.Profile, error) {

	profileFile, err := c.getProfileFilename()
	if err != nil {
		c.GetLogger().ErrorF("Failed to get profile filename: %s", err.Error())
		return nil, err
	}

	profileBody, err := os.ReadFile(profileFile)
	if err != nil {
		c.GetLogger().ErrorF("Failed to read profile file: %s", err.Error())
		return nil, err
	}

	profile, err := system.LoadJSONFromString[*bases.Profile](string(profileBody))
	if err != nil {
		c.GetLogger().ErrorF("Failed to unmarshal profile: %s", err.Error())
		return nil, err
	}

	if profile.Properties == nil {
		profile.Properties = make(map[string]string)
	}

	if profile.Secrets == nil {
		profile.Secrets = make(map[string]string)
	}

	return profile, nil

}

func (c *ProfileService) writeProfile(profile *bases.Profile) error {

	profileFile, err := c.getProfileFilename()
	if err != nil {
		c.GetLogger().ErrorF("Failed to get profile filename: %s", err.Error())
		return err
	}

	profileBody, err := system.DumpToJsonString(profile)
	if err != nil {
		c.GetLogger().ErrorF("Failed to marshal profile: %s", err.Error())
		return err
	}

	err = os.WriteFile(profileFile, []byte(profileBody), 0644)
	if err != nil {
		c.GetLogger().ErrorF("Failed to write profile file: %s", err.Error())
		return err
	}

	return nil

}
