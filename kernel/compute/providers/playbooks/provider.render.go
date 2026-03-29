package playbooks

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"path/filepath"
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

		if definition.Provision != "" || len(definition.Files) > 0 || definition.RawFiles != nil {
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

		if definition.Teardown != "" || len(definition.Files) > 0 || definition.RawFiles != nil {
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

	if verb == computeTypes.RunnerVerbProvision && (len(definition.Files) > 0 || definition.RawFiles != nil) {
		filesTasks, err := p.renderPlaybookFilesProvision(definition)
		if err != nil {
			return "", err
		}
		fullPlaybookTasks += filesTasks
	}

	fullPlaybookTasks += script

	if verb == computeTypes.RunnerVerbTeardown && (len(definition.Files) > 0 || definition.RawFiles != nil) {
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

var textFileExtensions = map[string]bool{
	".yaml": true, ".yml": true, ".json": true, ".txt": true,
	".html": true, ".htm": true, ".xml": true, ".css": true,
	".js": true, ".ts": true, ".sh": true, ".bash": true,
	".conf": true, ".cfg": true, ".ini": true, ".toml": true,
	".md": true, ".properties": true, ".env": true, ".sql": true,
	".py": true, ".go": true, ".java": true, ".rb": true,
	".php": true, ".pl": true, ".lua": true, ".j2": true,
	".jinja": true, ".jinja2": true, ".tpl": true, ".tmpl": true,
	".csv": true, ".tsv": true, ".log": true, ".service": true,
	".svg": true, ".tf": true, ".hcl": true,
}

func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		base := strings.ToLower(filepath.Base(path))
		switch base {
		case "dockerfile", "makefile", "vagrantfile", "rakefile", "gemfile", "procfile":
			return true
		}
		return false
	}
	return textFileExtensions[ext]
}

func (p *Provider) renderPlaybookFilesProvision(definition *types.PlaybookDefinition) (string, error) {
	paths := p.getPathsService()
	runDir := p.GetService().GetKernel().GetRunDirectory()
	baseDir := paths.CominePath(runDir, "var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name))
	filesDir := paths.CominePath(baseDir, "files")
	templatesDir := paths.CominePath(baseDir, "templates")

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

	tasks += fmt.Sprintf("- name: Create templates directory for %s\n", definition.Name)
	tasks += "  file:\n"
	tasks += fmt.Sprintf("    path: %s\n", templatesDir)
	tasks += "    state: directory\n"
	tasks += "    mode: \"0755\"\n"

	// Collect all subdirectories from both sources
	subdirs := make(map[string]bool)
	for relPath := range definition.Files {
		dir := paths.GetDirName(relPath)
		if dir != "." && dir != "" {
			subdirs[dir] = true
		}
	}
	if definition.RawFiles != nil {
		err := fs.WalkDir(definition.RawFiles, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			dir := paths.GetDirName(path)
			if dir != "." && dir != "" {
				subdirs[dir] = true
			}
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("failed to walk embedded files for %s: %w", definition.Name, err)
		}
	}

	for subdir := range subdirs {
		for _, parentDir := range []struct{ name, path string }{
			{"templates", templatesDir},
			{"files", filesDir},
		} {
			subdirPath := paths.CominePath(parentDir.path, subdir)
			tasks += fmt.Sprintf("- name: Create %s subdirectory %s for %s\n", parentDir.name, subdir, definition.Name)
			tasks += "  file:\n"
			tasks += fmt.Sprintf("    path: %s\n", subdirPath)
			tasks += "    state: directory\n"
			tasks += "    mode: \"0755\"\n"
		}
	}

	// Process PlaybookFileDefinition files (text with Jinja2 interpolation)
	for relPath, fileDef := range definition.Files {
		templatePath := paths.CominePath(templatesDir, relPath)
		targetPath := paths.CominePath(filesDir, relPath)
		mode := fileDef.Mode
		if mode == "" {
			mode = "0644"
		}

		tasks += fmt.Sprintf("- name: Write template source %s for %s\n", relPath, definition.Name)
		tasks += "  copy:\n"
		tasks += "    content: |\n"
		for _, line := range strings.Split(fileDef.Content, "\n") {
			tasks += fmt.Sprintf("      %s\n", line)
		}
		tasks += fmt.Sprintf("    dest: %s\n", templatePath)
		tasks += "    mode: \"0644\"\n"

		tasks += fmt.Sprintf("- name: Render file %s for %s\n", relPath, definition.Name)
		tasks += "  template:\n"
		tasks += fmt.Sprintf("    src: %s\n", templatePath)
		tasks += fmt.Sprintf("    dest: %s\n", targetPath)
		tasks += fmt.Sprintf("    mode: \"%s\"\n", mode)
	}

	// Process embed.FS raw files
	if definition.RawFiles != nil {
		err := fs.WalkDir(definition.RawFiles, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			data, readErr := definition.RawFiles.ReadFile(path)
			if readErr != nil {
				return readErr
			}

			targetPath := paths.CominePath(filesDir, path)

			if isTextFile(path) {
				// Text file: write to templates dir, then render with Jinja2
				templatePath := paths.CominePath(templatesDir, path)

				tasks += fmt.Sprintf("- name: Write template source %s for %s\n", path, definition.Name)
				tasks += "  copy:\n"
				tasks += "    content: |\n"
				for _, line := range strings.Split(string(data), "\n") {
					tasks += fmt.Sprintf("      %s\n", line)
				}
				tasks += fmt.Sprintf("    dest: %s\n", templatePath)
				tasks += "    mode: \"0644\"\n"

				tasks += fmt.Sprintf("- name: Render file %s for %s\n", path, definition.Name)
				tasks += "  template:\n"
				tasks += fmt.Sprintf("    src: %s\n", templatePath)
				tasks += fmt.Sprintf("    dest: %s\n", targetPath)
				tasks += "    mode: \"0644\"\n"
			} else {
				// Binary file: base64 encode, write encoded, decode with shell
				b64Content := base64.StdEncoding.EncodeToString(data)
				b64Path := targetPath + ".b64"

				tasks += fmt.Sprintf("- name: Write base64 encoded %s for %s\n", path, definition.Name)
				tasks += "  copy:\n"
				tasks += fmt.Sprintf("    content: \"%s\"\n", b64Content)
				tasks += fmt.Sprintf("    dest: %s\n", b64Path)
				tasks += "    mode: \"0644\"\n"

				tasks += fmt.Sprintf("- name: Decode binary file %s for %s\n", path, definition.Name)
				tasks += "  shell: |\n"
				tasks += fmt.Sprintf("    base64 -d < %s > %s\n", b64Path, targetPath)

				tasks += fmt.Sprintf("- name: Set permissions for %s for %s\n", path, definition.Name)
				tasks += "  file:\n"
				tasks += fmt.Sprintf("    path: %s\n", targetPath)
				tasks += "    mode: \"0644\"\n"

				tasks += fmt.Sprintf("- name: Remove encoded file %s for %s\n", path, definition.Name)
				tasks += "  file:\n"
				tasks += fmt.Sprintf("    path: %s\n", b64Path)
				tasks += "    state: absent\n"
			}

			return nil
		})
		if err != nil {
			return "", fmt.Errorf("failed to process embedded files for %s: %w", definition.Name, err)
		}
	}

	return tasks, nil
}

func (p *Provider) renderPlaybookFilesTeardown(definition *types.PlaybookDefinition) string {
	paths := p.getPathsService()
	runDir := p.GetService().GetKernel().GetRunDirectory()
	baseDir := paths.CominePath(runDir, "var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name))

	tasks := fmt.Sprintf("- name: Remove compute directory for %s\n", definition.Name)
	tasks += "  file:\n"
	tasks += fmt.Sprintf("    path: %s\n", baseDir)
	tasks += "    state: absent\n"

	return tasks
}
