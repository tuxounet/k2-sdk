package kubernetes

import (
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
)

func (p *Provider) Nuke() error {
	p.GetLogger().TraceF("Nuking up %s provider", ProviderKey)

	return nil
}

func (p *Provider) Setup() error {
	p.GetLogger().TraceF("Setting up %s provider", ProviderKey)

	return nil

}

func (p *Provider) Render() ([]computeTypes.RunnerDefinition, error) {
	p.GetLogger().TraceF("Rendering  %s provider", ProviderKey)

	return nil, nil
}
