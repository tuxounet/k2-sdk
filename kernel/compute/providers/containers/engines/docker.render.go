package engines

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
)

func (p *DockerEngine) RenderPlaybookTasks(definition types.ContainerDefinition, verb computeTypes.RunnerVerb) (string, error) {

	tasks := ""

	switch verb {
	case computeTypes.RunnerVerbProvision:
		files, err := p.renderProvisionContainerFilesTask(definition)
		if err != nil {
			return "", err
		}

		tasks = files

	case computeTypes.RunnerVerbStart:
		container, err := p.renderProvisionContainerTask(definition)
		if err != nil {
			return "", err
		}

		tasks = container

	case computeTypes.RunnerVerbStop:
		tasks = p.renderContainerDeclarationTask(definition, "stopped")

	case computeTypes.RunnerVerbTeardown:

		container := p.renderContainerDeclarationTask(definition, "absent")

		files, err := p.renderUnprovisionContainerFilesTask(definition)
		if err != nil {
			return "", err
		}

		tasks = files + container

	}

	return tasks, nil

}

func (p *DockerEngine) renderContainerDeclarationTask(definition types.ContainerDefinition, state string) string {

	task := fmt.Sprintf("- name: container %s\n", definition.Name)
	task += fmt.Sprintf("  community.docker.docker_container: #%s\n", definition.Name)
	task += fmt.Sprintf("    name: %d-%s\n", definition.Order, definition.Name)
	task += fmt.Sprintf("    state: %s\n", state)

	return task

}

func (p *DockerEngine) renderProvisionContainerTask(definition types.ContainerDefinition) (string, error) {
	localAddress, err := p.getLocalHostAddress()
	if err != nil {
		return "", err
	}

	task := fmt.Sprintf("- name: container %s\n", definition.Name)
	task += fmt.Sprintf("  community.docker.docker_container: #%s\n", definition.Name)
	task += fmt.Sprintf("    state: %s\n", "started")
	task += fmt.Sprintf("    name: %d-%s\n", definition.Order, definition.Name)
	task += fmt.Sprintf("    image: %s\n", definition.Image)
	task += fmt.Sprintf("    restart_policy: %s\n", "always")

	if definition.Capacities != nil {

		allowedCaps := []string{}
		droppedCaps := []string{}
		for _, cap := range *definition.Capacities {
			switch cap.Action {
			case "add":
				allowedCaps = append(allowedCaps, cap.Name)
			case "remove":
				droppedCaps = append(droppedCaps, cap.Name)
			}
		}

		if len(allowedCaps) > 0 {
			task += "    cap_add:\n"
			for _, allowed := range allowedCaps {
				task += fmt.Sprintf("      - %s\n", allowed)
			}
		}
		if len(droppedCaps) > 0 {
			task += "    cap_drop:\n"
			for _, dropped := range droppedCaps {
				task += fmt.Sprintf("      - %s\n", dropped)
			}
		}

	}

	if definition.Env != nil {
		task += "    env:\n"
		for key, value := range definition.Env {
			if strings.Contains(value, "\n") {
				task += fmt.Sprintf("      %s: |\n", key)
				for _, line := range strings.Split(value, "\n") {
					task += fmt.Sprintf("        %s\n", line)
				}
			} else {
				task += fmt.Sprintf("      %s: \"%s\"\n", key, value)
			}
		}
	}

	if definition.Command != nil && len(*definition.Command) > 0 {
		task += "    command:\n"
		for _, command := range *definition.Command {
			task += fmt.Sprintf("      - %s\n", command)
		}
	}

	if definition.Security != nil {
		if definition.Security.Privileged {
			task += "    privileged: yes\n"
		}

		if definition.Security.RunAsUser != "" {
			task += fmt.Sprintf("    user: %s\n", definition.Security.RunAsUser)
		}
	}

	if len(definition.Volumes) > 0 {

		task += "    volumes:\n"
		rootStore, err := p.getRootStore()
		if err != nil {
			return "", err
		}

		paths := p.getPathsService()

		for i, volume := range definition.Volumes {
			containerPath := volume.ContainerPath
			hostPath := ""

			switch volume.Binding.Type {
			case types.ContainerDefinitionVolumeBindingTypeContent:
				if volume.Binding.Content != "" {
					targetPath := paths.CominePath("var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name), "init", volume.ContainerPath)
					runDir := p.service.GetKernel().GetRunDirectory()
					hostPath = paths.CominePath(runDir, targetPath)
				} else {
					return "", fmt.Errorf("content value is not set for volume %s", volume.ContainerPath)
				}
			case types.ContainerDefinitionVolumeBindingTypeMount:

				if volume.Binding.HostPath != "" {
					hostPath = volume.Binding.HostPath
				} else {

					volumeName := volume.Name
					if volumeName == "" {
						volumeName = fmt.Sprintf("%s-%d", definition.Name, i)
					}
					//TODO: change var to data to get var non persistent
					targetPath := paths.CominePath("var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name), "volumes", volumeName)

					volumeFile := paths.CominePath(targetPath, ".volume")

					found, err := rootStore.Exists(volumeFile)
					if err != nil {
						return "", err
					}

					if !found {
						err = rootStore.WriteObject(paths.CominePath(targetPath, ".volume"), []byte(""))
						if err != nil {
							return "", err
						}
					}

					hostPath = paths.CominePath(p.service.GetKernel().GetRunDirectory(), targetPath)
				}

			case types.ContainerDefinitionVolumeBindingTypeEphemeral:

				volumeName := volume.Name
				if volumeName == "" {
					volumeName = fmt.Sprintf("%s-%d", definition.Name, i)
				}
				//TODO: change var to data to get var non persistent
				targetPath := paths.CominePath("tmp", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name), "volumes", volumeName)

				volumeFile := paths.CominePath(targetPath, ".volume")

				found, err := rootStore.Exists(volumeFile)
				if err != nil {
					return "", err
				}

				if !found {
					err = rootStore.WriteObject(paths.CominePath(targetPath, ".volume"), []byte(""))
					if err != nil {
						return "", err
					}
				}

				hostPath = paths.CominePath(p.service.GetKernel().GetRunDirectory(), targetPath)

			}

			if hostPath != "" {
				task += fmt.Sprintf("      - %s:%s\n", hostPath, containerPath)
			}
		}

	}

	if len(definition.Ports) > 0 || len(definition.Ingresses) > 0 {
		task += "    ports:\n"
		if len(definition.Ports) > 0 {

			for _, port := range definition.Ports {

				containerPort := port.ContainerPort
				containerProtcol := "TCP"
				if port.Protocol != "" {
					containerProtcol = port.Protocol
				}
				hostAddress := localAddress
				if port.HostAddress != "" {
					hostAddress = port.HostAddress
				}

				hostPort := port.HostPort
				if hostPort != "" {
					task += fmt.Sprintf("      - \"%s:%s/%s\"\n", fmt.Sprintf("%s:%s", hostAddress, hostPort), containerPort, strings.ToLower(containerProtcol))

				}
			}
		}

		if len(definition.Ingresses) > 0 {

			for _, ing := range definition.Ingresses {

				containerPort := ing.ContainerPort
				containerProtocol := "TCP"
				hostAddress := localAddress

				localPort, err := p.allocateLocalPort(definition, containerPort)

				if err != nil {
					p.logger.ErrorF("Failed to register ingress: %s on definition %s %s", err, definition.Name, definition.Order)
					return "", err
				}

				if localPort > 0 {

					ingressDef := &ingressTypes.IngressDefinition{
						AccessPolicy:  ing.AccessPolicy,
						ServiceHost:   hostAddress,
						ServicePort:   localPort,
						IngressPath:   ing.Path,
						RewritePath:   ing.RewritePath,
						CustomHandler: ing.CustomHandler,
					}
					err = p.ingressRegistar(ingressDef)
					if err != nil {
						p.logger.ErrorF("Failed to register ingress: %s on definition %s %s", err, definition.Name, definition.Order)
						return "", err
					}

					task += fmt.Sprintf("      - %s:%d/%s\n", fmt.Sprintf("%s:%d", hostAddress, localPort), containerPort, strings.ToLower(containerProtocol))
				}
			}
		}
	}

	if definition.Networks != nil && len(*definition.Networks) > 0 {
		task += "    networks:\n"
		for _, network := range *definition.Networks {
			task += fmt.Sprintf("      - name: %s\n", network)
		}
	}

	return task, nil
}

func (p *DockerEngine) renderProvisionContainerFilesTask(definition types.ContainerDefinition) (string, error) {

	files := make(map[string]string)

	if len(definition.Volumes) > 0 {
		for _, volumes := range definition.Volumes {
			if volumes.Binding.Type == "content" {
				if volumes.Binding.Content != "" {
					files[volumes.ContainerPath] = volumes.Binding.Content
				}
			}
		}
	}

	if len(files) > 0 {
		paths := p.getPathsService()

		runDir := p.service.GetKernel().GetRunDirectory()
		initFolder := paths.CominePath(runDir, "var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name), "init")

		fileProvisionTasks := ""

		fileProvisionTasks += fmt.Sprintf("- name: Create init folder for %s\n", definition.Name)
		fileProvisionTasks += fmt.Sprintf("  file:%s\n", "")
		fileProvisionTasks += fmt.Sprintf("    path: %s\n", initFolder)
		fileProvisionTasks += fmt.Sprintf("    state: directory%s\n", "")

		for containerPath, file := range files {
			fileFolder := paths.GetDirName(containerPath)
			fileFolderPath := paths.CominePath(initFolder, fileFolder)
			targetFile := paths.CominePath(initFolder, containerPath)
			fileProvisionTasks += fmt.Sprintf("- name: Create folder for %s\n", containerPath)
			fileProvisionTasks += fmt.Sprintf("  file:%s\n", "")
			fileProvisionTasks += fmt.Sprintf("    path: %s\n", fileFolderPath)
			fileProvisionTasks += fmt.Sprintf("    state: directory%s\n", "")
			fileProvisionTasks += fmt.Sprintf("- name: Writing init file for %s\n", containerPath)
			fileProvisionTasks += fmt.Sprintf("  copy:\n%s", "")
			fileProvisionTasks += fmt.Sprintf("    content: %s|\n", "")
			for _, line := range strings.Split(file, "\n") {
				fileProvisionTasks += fmt.Sprintf("     %s\n", line)
			}

			fileProvisionTasks += fmt.Sprintf("    dest: %s\n", targetFile)

		}

		return fileProvisionTasks, nil
	}

	return "", nil

}

func (p *DockerEngine) renderUnprovisionContainerFilesTask(definition types.ContainerDefinition) (string, error) {

	files := make([]string, 0)

	if len(definition.Volumes) > 0 {
		for _, volumes := range definition.Volumes {
			if volumes.Binding.Type == "content" {
				if volumes.Binding.Content != "" {
					files = append(files, volumes.ContainerPath)
				}
			}
		}
	}

	if len(files) > 0 {
		paths := p.getPathsService()

		runDir := p.service.GetKernel().GetRunDirectory()
		initFolder := paths.CominePath(runDir, "var", "compute", fmt.Sprintf("%d_%s", definition.Order, definition.Name), "init")

		fileProvisionTasks := fmt.Sprintf("- name: delete init folder for %s\n", definition.Name)
		fileProvisionTasks += fmt.Sprintf("  file:%s\n", "")
		fileProvisionTasks += fmt.Sprintf("    path: %s\n", initFolder)
		fileProvisionTasks += fmt.Sprintf("    state: absent%s\n", "")

		return fileProvisionTasks, nil
	}

	return "", nil

}

func (p *DockerEngine) allocateLocalPort(definition types.ContainerDefinition, port int) (int, error) {

	records, err := p.portMapStore.GetValue()
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
		return -1, errors.New("no more ports available")
	}

	record := types.PortsMapRecord{
		LocalPort:     localPort,
		ContainerName: definition.Name,
		ContainerPort: port,
	}

	allRecords = append(allRecords, record)

	err = p.portMapStore.SetValue(allRecords)
	if err != nil {
		p.logger.ErrorF("Failed to write portmaps: %s", err)
		return -1, err
	}

	return record.LocalPort, nil
}
