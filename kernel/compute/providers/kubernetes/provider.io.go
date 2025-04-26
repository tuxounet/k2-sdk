package kubernetes

import (
	"path/filepath"

	"github.com/tuxounet/k2-sdk/kernel/config"
)

func (p *Provider) getConfigService() *config.Service {
	return p.GetService().GetKernel().GetService(config.ServiceKey).(*config.Service)
}

func (p *Provider) getIsEmbeddedEnabled() bool {
	val, err := p.getConfigService().GetAsBool("host.compute.kubernetes.useEmbedded")
	if err != nil {
		p.GetLogger().ErrorF("unable to get embedded kubernetes config value: %s", err)
		return false
	}
	return val
}
func (p *Provider) getKubeConfigValue() string {
	defaultValue := filepath.Join(p.GetService().GetKernel().GetRunDirectory(), "home", ".kube", "config")
	value, err := p.getConfigService().GetAsStringOrDefault("host.compute.kubernetes.kubeConfig", defaultValue)
	if err != nil {
		p.GetLogger().WarnF("unable to get kubeConfig config value: %s, using default %s", err, defaultValue)
		return defaultValue

	}

	return value
}

func (p *Provider) getKubeApiPort() int {
	defaultValue := 6443
	value, err := p.getConfigService().GetAsInt("host.compute.kubernetes.kubeApiPort")
	if err != nil {
		p.GetLogger().WarnF("unable to get kubeApiPort config value: %s, using default %d", err, defaultValue)
		return defaultValue
	}

	return value
}
