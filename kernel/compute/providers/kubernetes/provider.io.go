package kubernetes

import (
	"path/filepath"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
	"github.com/tuxounet/k2-sdk/kernel/config"
	storesTypes "github.com/tuxounet/k2-sdk/kernel/storage/stores/types"
)

func (s *Provider) getPortsForwardsStore() storesTypes.IBaseObjectStore[[]types.PortsForwardRecord] {
	return s.GetData("forwards").(storesTypes.IBaseObjectStore[[]types.PortsForwardRecord])
}

func (s *Provider) getLocalHostAddress() (string, error) {
	defaultValue := "127.0.0.1"
	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	return configService.GetAsStringOrDefault("host.address", defaultValue)

}

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

	if value == "" {
		p.GetLogger().WarnF("kubeConfig config value is empty, using default %s", defaultValue)
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
	value, err := p.getConfigService().GetAsInt("host.compute.kubernetes.embedded.ports.api")
	if err != nil {
		p.GetLogger().WarnF("unable to get kubeApiPort config value: %s, using default %d", err, defaultValue)
		return defaultValue
	}

	return value
}

func (p *Provider) getKubeIngressPortPlain() int {
	defaultValue := 80
	value, err := p.getConfigService().GetAsInt("host.compute.kubernetes.embedded.ports.http")
	if err != nil {
		p.GetLogger().WarnF("unable to get kubeIngressPortPlain config value: %s, using default %d", err, defaultValue)
		return defaultValue
	}

	return value
}
func (p *Provider) getKubeIngressPortTls() int {
	defaultValue := 443
	value, err := p.getConfigService().GetAsInt("host.compute.kubernetes.embedded.ports.https")
	if err != nil {
		p.GetLogger().WarnF("unable to get kubeIngressPortTls config value: %s, using default %d", err, defaultValue)
		return defaultValue
	}

	return value
}

func (s *Provider) getHostPortStart() (int, error) {

	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.kubernetes.port.start")
	if err != nil {
		return -1, err
	}
	return port, nil

}
func (s *Provider) getHostPortEnd() (int, error) {
	kernel := s.GetService().GetKernel()
	configService := kernel.GetService(config.ServiceKey).(*config.Service)

	port, err := configService.GetAsInt("host.compute.kubernetes.port.end")
	if err != nil {
		return -1, err
	}
	return port, nil
}

func (s *Provider) getForwarders() []*types.PortForwarder {
	return s.GetData("forwarders").([]*types.PortForwarder)
}

func (s *Provider) setForwarders(forwarders []*types.PortForwarder) {
	s.SetData("forwarders", forwarders)
}
