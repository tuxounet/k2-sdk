package secrets

import "github.com/tuxounet/k2-sdk/kernel/secrets/bases"

func (c *SecretsService) Init() error {
	profiles := c.getProfileService()

	found, err := profiles.HasSecret("ssh.private")
	if err != nil {
		c.GetLogger().ErrorF("Failed to check for SSH private key: %s", err.Error())
		return err
	}

	if !found {
		c.GetLogger().Info("Generating new SSH key pair")
		kp, err := bases.NewKeyPair()
		if err != nil {
			c.GetLogger().ErrorF("Failed to generate SSH key pair: %s", err.Error())
			return err
		}
		err = c.SetSecret("ssh.private", kp.PrivateKey)
		if err != nil {
			c.GetLogger().ErrorF("Failed to set ssh.private secret: %s", err.Error())
			return err
		}

		err = c.SetSecret("ssh.public", kp.PublicKey)
		if err != nil {
			c.GetLogger().ErrorF("Failed to set ssh.public secret: %s", err.Error())
			return err
		}

	}

	return nil
}
