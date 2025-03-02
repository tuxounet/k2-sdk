package secrets

import "github.com/tuxounet/k2-sdk/kernel/profile"

func (c *SecretsService) getProfileService() *profile.ProfileService {
	profiles := c.GetKernel().GetService(profile.ServiceKey)
	if profiles == nil {
		c.GetLogger().Warn("Profile service not found")
		return nil
	}

	profileService := profiles.(*profile.ProfileService)

	return profileService
}
