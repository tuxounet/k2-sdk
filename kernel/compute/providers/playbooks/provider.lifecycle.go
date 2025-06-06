package playbooks

import (
	"github.com/tuxounet/k2-sdk/system"
)

func (p *Provider) Setup() error {
	p.GetLogger().TraceF("Setting up %s provider", ProviderKey)

	err := p.Nuke()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke definitions store: %s", err)
		return err
	}

	checkCmd := system.NewCmdCall(p.GetLogger(), "ansible-playbook", "--version")

	exit, err := system.OsExecWithExitCode(checkCmd)
	if err != nil {
		p.GetLogger().ErrorF("%s provider isReady failed: %s", ProviderKey, err)
		return err
	}

	if exit != 0 {
		p.GetLogger().ErrorF("%s provider isReady failed: exit code %d", ProviderKey, exit)
		return nil
	}

	p.GetLogger().DebugF("%s provider setup done", ProviderKey)
	return nil

}
