package kubernetes

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
	"github.com/tuxounet/k2-sdk/system"
)

//go:embed verbs/setup.yaml
var setupPlaybook string

//go:embed verbs/nuke.yaml
var nukePlaybook string

//go:embed verbs/start.yaml
var startPlaybook string

//go:embed verbs/stop.yaml
var stopPlaybook string

func (p *Provider) Render() ([]computeTypes.RunnerDefinition, error) {

	definitions := p.GetDefinitions()

	if len(definitions) == 0 {
		p.GetLogger().DebugF("No definitions found for %s provider", ProviderKey)
		return nil, nil
	}
	runners := make([]computeTypes.RunnerDefinition, 0)

	if p.getIsEmbeddedEnabled() {

		setupScript, err := p.renderPlaybookProvider(setupPlaybook)
		if err != nil {
			return nil, err
		}
		nukeScript, err := p.renderPlaybookProvider(nukePlaybook)
		if err != nil {
			return nil, err
		}

		setupRunner := computeTypes.RunnerDefinition{
			Order:     0,
			Name:      "kube",
			Provider:  ProviderKey,
			Provision: setupScript,
			Start:     "",
			Stop:      "",
			Teardown:  nukeScript,
		}
		runners = append(runners, setupRunner)

	}
	configService := p.getConfigService()

	for _, definition := range definitions {
		newRunnerDefinition := computeTypes.RunnerDefinition{
			Order:     definition.Order,
			Name:      definition.Name,
			Provider:  ProviderKey,
			Provision: "",
			Start:     "",
			Stop:      "",
			Teardown:  "",
		}

		startScript, err := p.renderPlaybookTasks(&definition, startPlaybook)
		if err != nil {
			return nil, err
		}
		newRunnerDefinition.Start = startScript
		stopScript, err := p.renderPlaybookTasks(&definition, stopPlaybook)
		if err != nil {
			return nil, err
		}

		newRunnerDefinition.Stop = stopScript

		if len(definition.Ingresses) > 0 {
			localAddress, err := p.getLocalHostAddress()
			if err != nil {
				return nil, err
			}

			ingressRegistar := p.GetIngressRegistar()
			for _, ing := range definition.Ingresses {
				if ing.ServiceNamespace == "" {
					ing.ServiceNamespace, err = configService.GetAsStringOrDefault("host.compute.kubernetes.ingress.serviceNamespace", "kube-system")
					if err != nil {
						return nil, fmt.Errorf("failed to get ingress service namespace: %s", err)
					}
				}
				if ing.ServiceName == "" {
					ing.ServiceName, err = configService.GetAsStringOrDefault("host.compute.kubernetes.ingress.serviceName", "traefik")
					if err != nil {
						return nil, fmt.Errorf("failed to get ingress service name: %s", err)
					}
				}
				if ing.ServicePort == 0 {
					ing.ServicePort, err = configService.GetAsIntOrDefault("host.compute.kubernetes.ingress.servicePort", 8000)
					if err != nil {
						return nil, fmt.Errorf("failed to get ingress service port: %s", err)
					}
				}

				localPort, err := p.allocateIngressPort(ing.IngressPath, ing.ServiceNamespace, ing.ServiceName, ing.ServicePort)
				if err != nil {
					return nil, fmt.Errorf("failed to allocate local port: %s", err)
				}

				customHandler, err := p.getIngressHandler(localPort)
				if err != nil {
					return nil, fmt.Errorf("failed to get ingress handler: %s", err)
				}

				ingressDef := &ingressTypes.IngressDefinition{
					AccessPolicy: ing.AccessPolicy,
					IngressPath:  ing.IngressPath,
					ServiceHost:  localAddress,
					ServicePort:  localPort,
					CustomHandler: func(_ ingressTypes.IngressDefinition) gin.HandlerFunc {
						return customHandler
					},
				}

				err = ingressRegistar(ingressDef)
				if err != nil {
					return nil, fmt.Errorf("failed to get ingress port: %s", err)
				}
			}

		}

		if len(definition.Ports) > 0 {
			for _, port := range definition.Ports {
				err := p.allocateLocalPort(port.LocalPort, port.ServiceNamespace, port.ServiceName, port.ServicePort)
				if err != nil {
					return nil, fmt.Errorf("failed to allocate local port: %s", err)
				}

			}
		}

		runners = append(runners, newRunnerDefinition)
	}

	return runners, nil
}

func (p *Provider) getTemplateValues() map[string]any {
	return map[string]any{
		"kubecontext":  strings.ToLower(p.GetService().GetKernel().GetApp().GetName()),
		"kubeconfig":   p.getKubeConfigValue(),
		"kubeNetworks": p.getKubeNetworks(),
	}
}

func (p *Provider) renderPlaybookProvider(script string) (string, error) {
	fullPlaybookProvider := fmt.Sprintf("# [%s] Ansible Playbook\n", ProviderKey)

	values := p.getTemplateValues()

	untemplated, err := system.UnTemplateWithGoTemplate(script, values)
	if err != nil {
		return "", fmt.Errorf("failed to untemplate script: %s", err)
	}
	fullPlaybookProvider += untemplated

	return fullPlaybookProvider, nil
}

func (p *Provider) renderPlaybookTasks(definition *types.NamespaceDefinition, script string) (string, error) {
	fullPlaybookTasks := fmt.Sprintf("# [%s] Ansible Playbook %d of %s\n", ProviderKey, definition.Order, definition.Name)

	values := p.getTemplateValues()
	values["namespace"] = definition.Name

	entries, err := walkTemplates(definition.Templates)
	if err != nil {
		return "", fmt.Errorf("failed to walk templates: %s", err)
	}
	values["templates"] = entries

	untemplated, err := system.UnTemplateWithGoTemplate(script, values)
	if err != nil {
		return "", fmt.Errorf("failed to untemplate script: %s", err)
	}
	fullPlaybookTasks += untemplated

	return fullPlaybookTasks, nil

}

func walkTemplates(fs *embed.FS) (map[string]any, error) {
	templates := make(map[string]any)

	var readDirRecursive func(string) error
	readDirRecursive = func(dir string) error {
		entries, err := fs.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fullPath := entry.Name()
			if entry.IsDir() {
				if err := readDirRecursive(filepath.Join(dir, fullPath)); err != nil {
					return err
				}
			} else {
				content, err := fs.ReadFile(filepath.Join(dir, fullPath))
				if err != nil {
					continue
				}
				tpl, err := system.LoadYamlFromString[any](string(content))
				if err != nil {
					return fmt.Errorf("failed to load yaml from string: %s", err)
				}

				templates[filepath.Join(dir, fullPath)] = tpl
			}
		}
		return nil
	}

	if err := readDirRecursive("."); err != nil {
		return nil, fmt.Errorf("failed to read templates: %s", err)
	}

	return templates, nil

}

func (p *Provider) allocateIngressPort(ingressPath string, serviceNamespace string, serviceName string, servicePort int) (int, error) {

	records, err := p.getPortsForwardsStore().GetValue()
	if err != nil {
		return -1, err
	}
	allRecords := *records

	portStart, err := p.getHostPortStart()
	if err != nil {
		return -1, err
	}
	index := len(allRecords)

	localPort := portStart + index

	portEnd, err := p.getHostPortEnd()
	if err != nil {
		return -1, err
	}

	if localPort > portEnd {
		return -1, fmt.Errorf("no more ports available")
	}

	record := types.PortsForwardRecord{
		Path:             ingressPath,
		ServiceNamespace: serviceNamespace,
		ServiceName:      serviceName,
		ServicePort:      servicePort,
		LocalPort:        localPort,
	}

	allRecords = append(allRecords, record)

	err = p.getPortsForwardsStore().SetValue(allRecords)
	if err != nil {
		p.GetLogger().ErrorF("Failed to write portmaps: %s", err)
		return -1, err
	}

	return record.LocalPort, nil
}

func (p *Provider) allocateLocalPort(localPort int, serviceNamespace string, serviceName string, servicePort int) error {

	records, err := p.getPortsForwardsStore().GetValue()
	if err != nil {
		return err
	}
	allRecords := *records

	found := false
	for _, record := range allRecords {
		if record.LocalPort == localPort {
			found = true
			break
		}
	}
	if found {
		return fmt.Errorf("port %d already allocated", localPort)
	}

	record := types.PortsForwardRecord{
		Path:             fmt.Sprintf(":%d", localPort),
		ServiceNamespace: serviceNamespace,
		ServiceName:      serviceName,
		ServicePort:      servicePort,
		LocalPort:        localPort,
	}

	allRecords = append(allRecords, record)

	err = p.getPortsForwardsStore().SetValue(allRecords)
	if err != nil {
		p.GetLogger().ErrorF("Failed to write portmaps: %s", err)
		return err
	}

	return nil
}

func (p *Provider) getIngressHandler(localPort int) (gin.HandlerFunc, error) {

	forwards, err := p.getPortsForwardsStore().GetValue()
	if err != nil {
		return nil, err
	}

	var instance *types.PortsForwardRecord
	allForwards := *forwards
	for _, forward := range allForwards {
		if forward.LocalPort == localPort {
			instance = &forward
			break
		}
	}
	if instance == nil {
		return nil, fmt.Errorf("no port forward found for local port %d", localPort)
	}

	handler := func(c *gin.Context) {
		forwarders := p.getForwarders()
		for _, forwarder := range forwarders {
			if forwarder.Record.ServiceNamespace == instance.ServiceNamespace &&
				forwarder.Record.ServiceName == instance.ServiceName &&
				forwarder.Record.ServicePort == instance.ServicePort {
				p.GetLogger().DebugF("Forwarding request to %s:%d", instance.ServiceName, instance.ServicePort)
				if !forwarder.IsReady() {
					err = forwarder.Mount()
					if err != nil {
						p.GetLogger().ErrorF("Failed to mount forwarder: %s", err)
						c.Status(http.StatusBadGateway)
						return
					}
				}

				err = forwarder.ForwardRequest(c)
				if err != nil {
					p.GetLogger().ErrorF("Failed to forward request: %s", err)
					c.Status(500)
					return
				}
				return
			}
		}

		p.GetLogger().WarnF("No forwarder found for %s:%d", instance.ServiceName, instance.ServicePort)
		c.Status(404)
	}

	return handler, nil
}
