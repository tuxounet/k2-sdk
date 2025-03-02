package containers

func (p *Provider) Nuke() error {
	p.GetLogger().TraceF("Nuking up %s provider", ProviderKey)

	p.killRootlessPort()

	err := p.getPortsMapSore().Nuke()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke portsmap store: %s", err)
		return err
	}
	return nil
}

func (p *Provider) Setup() error {
	p.GetLogger().TraceF("Setting up %s provider", ProviderKey)

	_, err := p.listContainers()
	if err != nil {
		p.GetLogger().ErrorF("%s provider is not ready : %s", ProviderKey, err.Error())
		p.GetLogger().InfoF("Try to install podman with the following command: 'sudo apt install podman podman-compose'")
		return err
	}

	p.GetLogger().DebugF("%s provider setup done", ProviderKey)
	return nil

}
