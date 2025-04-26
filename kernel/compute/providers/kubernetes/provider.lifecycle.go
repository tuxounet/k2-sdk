package kubernetes

func (p *Provider) Nuke() error {
	p.GetLogger().TraceF("Nuking up %s provider", ProviderKey)

	return nil
}

func (p *Provider) Setup() error {
	p.GetLogger().TraceF("Setting up %s provider", ProviderKey)

	return nil

}
