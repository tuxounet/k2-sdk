package engines

import (
	"github.com/tuxounet/k2-sdk/system"
)

func (e *PodmanEngine) withPodmanCliCmd(args ...string) *system.CmdCall {
	log := e.logger.CreateSubLogger("exec")
	cmdCall := system.NewCmdCall(log, "podman", args...)

	runDir := e.service.GetKernel().GetRunDirectory()
	cmdCall.Cwd = &runDir
	return cmdCall
}

func (e *PodmanEngine) listContainers() ([]*listPodmanContainerOutput, error) {

	out, err := system.JsonCommandOutput[[]map[string]any](e.withPodmanCliCmd("ps", "--all", "--format", "json"))
	if err != nil {
		e.logger.ErrorF("%s provider listContainers failed: %s", e.Name, err)
		return nil, err
	}

	ret := []*listPodmanContainerOutput{}

	for _, container := range out {
		ret = append(ret, &listPodmanContainerOutput{
			ID:    container["Id"].(string),
			Image: container["Image"].(string),
			State: container["State"].(string),
			Names: container["Names"].([]any),
		})

	}

	return ret, nil

}

type listPodmanContainerOutput struct {
	ID    string        `json:"ID"`
	Image string        `json:"Image"`
	State string        `json:"State"`
	Names []interface{} `json:"Names"`
}

func (e *PodmanEngine) killRootlessPort() {
	processName := "rootlessport"
	e.logger.TraceF("trying killing process: %s", processName)
	killCmd := system.NewCmdCall(e.logger, "pkill", "-f", processName)

	err := system.OsExec(killCmd)
	if err != nil {
		e.logger.WarnF("Failed to kill process %s: %s", processName, err)
	}

}
