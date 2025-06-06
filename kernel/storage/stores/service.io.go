package stores

import (
	"github.com/tuxounet/k2-sdk/kernel/profile"
)

func (c *Service) getProfileService() *profile.ProfileService {
	profileService := c.GetKernel().GetService(profile.ServiceKey)
	if profileService == nil {
		c.GetLogger().Warn("Profile service not found")
		return nil
	}

	return profileService.(*profile.ProfileService)
}
