package secrets

import (
	"fmt"

	"github.com/tuxounet/k2-sdk/kernel/secrets/bases"
)

func (c *SecretsService) Init() error {
	profiles := c.getProfileService()

	found, err := profiles.HasSecret("ssh.private")
	if err != nil {
		c.GetLogger().ErrorF("failed to check for SSH private key: %s", err.Error())
		return fmt.Errorf("secrets: failed to check SSH private key existence: %w", err)
	}

	if !found {
		c.GetLogger().Info("Generating new SSH key pair")
		kp, err := bases.NewKeyPair()
		if err != nil {
			c.GetLogger().ErrorF("failed to generate SSH key pair: %s", err.Error())
			return fmt.Errorf("secrets: failed to generate SSH key pair: %w", err)
		}
		err = c.SetSecret("ssh.private", kp.PrivateKey)
		if err != nil {
			c.GetLogger().ErrorF("failed to set ssh.private secret: %s", err.Error())
			return fmt.Errorf("secrets: failed to persist ssh.private secret: %w", err)
		}

		err = c.SetSecret("ssh.public", kp.PublicKey)
		if err != nil {
			c.GetLogger().ErrorF("failed to set ssh.public secret: %s", err.Error())
			return fmt.Errorf("secrets: failed to persist ssh.public secret: %w", err)
		}

	}

	return nil
}
