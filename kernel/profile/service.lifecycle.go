package profile

import (
	"fmt"
	"os"

	"github.com/tuxounet/k2-sdk/kernel/profile/bases"
)

func (c *ProfileService) Init() error {
	profile, err := c.readProfile()
	if err != nil {
		if os.IsNotExist(err) {
			profile = &bases.Profile{
				Name:       c.GetKernel().GetApp().GetName(),
				Version:    c.GetKernel().GetApp().GetVersion(),
				Properties: make(map[string]string),
				Secrets:    make(map[string]string),
			}

			err = c.writeProfile(profile)
			if err != nil {
				c.GetLogger().ErrorF("failed to write initial profile: %s", err.Error())
				return fmt.Errorf("profile: failed to write initial profile: %w", err)
			}
		} else {
			c.GetLogger().ErrorF("failed to read profile: %s", err.Error())
			return fmt.Errorf("profile: failed to read profile: %w", err)
		}
	}

	profile.Name = c.GetKernel().GetApp().GetName()
	profile.Version = c.GetKernel().GetApp().GetVersion()

	err = c.writeProfile(profile)
	if err != nil {
		c.GetLogger().ErrorF("failed to write profile: %s", err.Error())
		return fmt.Errorf("profile: failed to write profile: %w", err)
	}

	return nil
}
