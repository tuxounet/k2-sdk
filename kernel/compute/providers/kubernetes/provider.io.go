package kubernetes

import (
	"path/filepath"

	"github.com/tuxounet/k2-sdk/kernel/config"
)

func (p *Provider) getConfigService() *config.Service {
	return p.GetService().GetKernel().GetService(config.ServiceKey).(*config.Service)
}

func (p *Provider) getIsEmbeddedEnabled() bool {
	val, err := p.getConfigService().GetAsBool("host.compute.kubernetes.embedded.enabled")
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

func (p *Provider) getKubeImage() string {
	defaultValue := "docker.io/rancher/k3s:v1.30.11-k3s1"
	value, err := p.getConfigService().GetAsStringOrDefault("host.compute.kubernetes.embedded.image", defaultValue)
	if err != nil {
		p.GetLogger().WarnF("unable to get kube image config value: %s, using default %s", err, defaultValue)
		return defaultValue

	}

	return value
}

func (p *Provider) getKubeNetworks() []any {
	defaultValue := []any{
		map[string]string{
			"subnet":  "172.99.27.0/24",
			"gateway": "172.99.27.2",
			"iprange": "172.99.27.0/26",
		},
	}

	value := p.getConfigService().Get("host.compute.kubernetes.embedded.networks")
	if value == nil {
		p.GetLogger().WarnF("unable to get kube network config value, using default %v", defaultValue)
		return defaultValue
	}
	if _, ok := value.([]any); !ok {
		p.GetLogger().WarnF("kube network config value is not a list, using default %v", defaultValue)
		return defaultValue
	}
	if len(value.([]any)) == 0 {
		p.GetLogger().WarnF("kube network config value is empty, using default %v", defaultValue)
		return defaultValue
	}

	return value.([]any)
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
