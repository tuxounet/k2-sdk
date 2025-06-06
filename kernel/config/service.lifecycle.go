package config

func (c *Service) Init() error {

	err := c.initDefaultConfig()
	if err != nil {
		c.GetLogger().ErrorF("Error initializing default config: %s", err.Error())
		return err
	}

	return nil
}
