package playbooks

import (
	"fmt"
	"strings"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
)

func (p *Provider) Render() ([]computeTypes.RunnerDefinition, error) {

	p.GetLogger().DebugF("[RENDER] Rendering %s provider", ProviderKey)

	definitions := p.GetDefinitions()
	runners := make([]computeTypes.RunnerDefinition, 0)

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

		if definition.Provision != "" || len(definition.Files) > 0 {
			provisionScript := definition.Provision
			provisionTasks, err := p.renderPlaybookTasks(&definition, computeTypes.RunnerVerbProvision, provisionScript)
			if err != nil {
				return nil, err
			}
			newRunnerDefinition.Provision = provisionTasks
		}

		if definition.Start != "" {
			startTasks, err := p.renderPlaybookTasks(&definition, computeTypes.RunnerVerbStart, definition.Start)
			if err != nil {
				return nil, err
			}
			newRunnerDefinition.Start = startTasks
		}

		if definition.Stop != "" {
			stopTasks, err := p.renderPlaybookTasks(&definition, computeTypes.RunnerVerbStop, definition.Stop)
			if err != nil {
				return nil, err
			}
			newRunnerDefinition.Stop = stopTasks
		}

		if definition.Teardown != "" || len(definition.Files) > 0 {
			teardownScript := definition.Teardown
			teardownTasks, err := p.renderPlaybookTasks(&definition, computeTypes.RunnerVerbTeardown, teardownScript)
			if err != nil {
				return nil, err
			}
			newRunnerDefinition.Teardown = teardownTasks

		}
		runners = append(runners, newRunnerDefinition)
	}

	p.GetLogger().DebugF("Rendered %d runners for %s provider", len(runners), ProviderKey)

	return runners, nil
}

func (p *Provider) renderPlaybookTasks(definition *types.PlaybookDefinition, verb computeTypes.RunnerVerb, script string) (string, error) {
	fullPlaybookTasks := fmt.Sprintf("# [%d] Ansible Playbook %s of %s\n", definition.Order, verb, definition.Name)

	if verb == computeTypes.RunnerVerbProvision && len(definition.Files) > 0 {
		filesTasks, err := p.renderPlaybookFilesProvision(definition)
		if err != nil {
			return "", err
		}
		fullPlaybookTasks += filesTasks
	}

	fullPlaybookTasks += script

	if verb == computeTypes.RunnerVerbTeardown && len(definition.Files) > 0 {
		cleanupTasks := p.renderPlaybookFilesTeardown(definition)
		fullPlaybookTasks += "\n" + cleanupTasks
	}

	return fullPlaybookTasks, nil

}

func (p *Provider) sanitizeVarName(name string) string {
	result := ""
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			result += string(c)
		} else {
			result += "_"
		}
	}
	return result
}

func (p *Provider) renderPlaybookFilesProvision(definition *types.PlaybookDefinition) (string, error) {
	paths := p.getPathsService()
	runDir := p.GetService().GetKernel().GetRunDirectory()
	filesDir := paths.CominePath(runDir, "var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name), "files")

	varName := fmt.Sprintf("playbook_files_%s", p.sanitizeVarName(definition.Name))

	tasks := ""

	tasks += fmt.Sprintf("- name: Set files path for %s\n", definition.Name)
	tasks += "  set_fact:\n"
	tasks += fmt.Sprintf("    %s: %s\n", varName, filesDir)

	tasks += fmt.Sprintf("- name: Create files directory for %s\n", definition.Name)
	tasks += "  file:\n"
	tasks += fmt.Sprintf("    path: %s\n", filesDir)
	tasks += "    state: directory\n"
	tasks += "    mode: \"0755\"\n"

	subdirs := make(map[string]bool)
	for relPath := range definition.Files {
		dir := paths.GetDirName(relPath)
		if dir != "." && dir != "" {
			subdirs[dir] = true
		}
	}

	for subdir := range subdirs {
		subdirPath := paths.CominePath(filesDir, subdir)
		tasks += fmt.Sprintf("- name: Create subdirectory %s for %s\n", subdir, definition.Name)
		tasks += "  file:\n"
		tasks += fmt.Sprintf("    path: %s\n", subdirPath)
		tasks += "    state: directory\n"
		tasks += "    mode: \"0755\"\n"
	}

	for relPath, fileDef := range definition.Files {
		targetPath := paths.CominePath(filesDir, relPath)
		mode := fileDef.Mode
		if mode == "" {
			mode = "0644"
		}
		tasks += fmt.Sprintf("- name: Write file %s for %s\n", relPath, definition.Name)
		tasks += "  copy:\n"
		tasks += "    content: |\n"
		for _, line := range strings.Split(fileDef.Content, "\n") {
			tasks += fmt.Sprintf("      %s\n", line)
		}
		tasks += fmt.Sprintf("    dest: %s\n", targetPath)
		tasks += fmt.Sprintf("    mode: \"%s\"\n", mode)
	}

	return tasks, nil
}

func (p *Provider) renderPlaybookFilesTeardown(definition *types.PlaybookDefinition) string {
	paths := p.getPathsService()
	runDir := p.GetService().GetKernel().GetRunDirectory()
	filesDir := paths.CominePath(runDir, "var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name), "files")

	tasks := fmt.Sprintf("- name: Remove files directory for %s\n", definition.Name)
	tasks += "  file:\n"
	tasks += fmt.Sprintf("    path: %s\n", filesDir)
	tasks += "    state: absent\n"

	return tasks
}
