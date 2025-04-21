package containers

func (p *Provider) Nuke() error {
	p.GetLogger().TraceF("Nuking up %s provider", ProviderKey)

	engine := p.getContainerEngine()
	err := engine.Nuke()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke %s provider: %s", ProviderKey, err)
		return err
	}

	err = p.getPortsMapSore().Nuke()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke portsmap store: %s", err)
		return err
	}
	return nil
}

func (p *Provider) Setup() error {
	p.GetLogger().TraceF("Setting up %s provider", ProviderKey)

	engine := p.getContainerEngine()
	err := engine.Setup()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke %s provider: %s", ProviderKey, err)
		return err
	}

	p.GetLogger().DebugF("%s provider setup done", ProviderKey)
	return nil

}
