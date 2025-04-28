package kubernetes

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/kubernetes/types"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
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

		runners = append(runners, newRunnerDefinition)
	}

	return runners, nil
}

func (p *Provider) getTemplateValues() map[string]any {
	return map[string]any{
		"kubecontext":  strings.ToLower(p.GetService().GetKernel().GetApp().GetName()),
		"kubeconfig":   p.getKubeConfigValue(),
		"kubeApiPort":  fmt.Sprintf("%d", p.getKubeApiPort()),
		"kubeImage":    p.getKubeImage(),
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
