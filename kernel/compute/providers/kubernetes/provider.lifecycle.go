package kubernetes

import (
	"strings"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
)

func (p *Provider) Init() error {
	p.GetLogger().TraceF("Initializing %s provider", ProviderKey)
	err := p.getPortsForwardsStore().Nuke()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke ports forwards store: %s", err)
		return err
	}
	return nil
}

func (p *Provider) Nuke() error {
	p.GetLogger().TraceF("Nuking up %s provider", ProviderKey)

	return nil
}

func (p *Provider) Setup() error {
	p.GetLogger().TraceF("Setting up %s provider", ProviderKey)

	return nil

}

func (p *Provider) Start() error {
	p.GetLogger().TraceF("Starting %s provider", ProviderKey)

	forwards, err := p.getPortsForwardsStore().GetValue()
	if err != nil {
		p.GetLogger().ErrorF("Failed to get ports forwards store: %s", err)
		return err
	}
	hostAddress, err := p.getLocalHostAddress()
	if err != nil {
		p.GetLogger().ErrorF("Failed to get host address: %s", err)
		return err
	}
	forwarders := make([]*types.PortForwarder, 0)
	for _, forward := range *forwards {

		p.GetLogger().ErrorF("Starting port forward %v", forward)
		forwarder := types.NewPortForwarder(forward, p.getKubeConfigValue(), p.GetLogger(), hostAddress)
		forwarders = append(forwarders, forwarder)

		if strings.HasPrefix(forwarder.Record.Path, ":") {
			err := forwarder.Mount()
			if err != nil {
				p.GetLogger().ErrorF("Failed to mount port forward %v: %s", forward, err)
				return err
			}

		}

	}

	p.setForwarders(forwarders)
	return nil
}

func (p *Provider) Stop() error {
	p.GetLogger().TraceF("Stopping %s provider", ProviderKey)

	for _, forward := range p.getForwarders() {
		p.GetLogger().ErrorF("Stoping port forward %v", forward)
		err := forward.Stop()
		if err != nil {
			p.GetLogger().ErrorF("Failed to stop port forward %v: %s", forward, err)
			return err
		}
	}

	return nil
}
