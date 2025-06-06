package engines

import (
	"strings"

	"github.com/tuxounet/k2-sdk/system"
)

func (e *DockerEngine) withPodmanCliCmd(args ...string) *system.CmdCall {
	log := e.logger.CreateSubLogger("exec")
	cmdCall := system.NewCmdCall(log, "docker", args...)

	runDir := e.service.GetKernel().GetRunDirectory()
	cmdCall.Cwd = &runDir
	return cmdCall
}

func (e *DockerEngine) listContainers() ([]*listDockerContainerOutput, error) {

	raw, err := system.RawCommandOutput(e.withPodmanCliCmd("ps", "--all", "--format", "json"))
	if err != nil {
		e.logger.ErrorF("%s provider listContainers failed: %s", e.Name, err)
		return nil, err
	}
	if raw == "" {
		return make([]*listDockerContainerOutput, 0), nil
	}

	lines := strings.Split(raw, "\n")
	if len(lines) == 0 {
		return make([]*listDockerContainerOutput, 0), nil
	}
	ret := []*listDockerContainerOutput{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		out, err := system.LoadJSONFromString[map[string]string](line)
		if err != nil {
			e.logger.ErrorF("%s provider listContainers failed: %s", e.Name, err)
			return nil, err
		}

		ret = append(ret, &listDockerContainerOutput{
			ID:    out["ID"],
			Image: out["Image"],
			State: out["State"],
			Names: out["Names"],
		})

	}
	return ret, nil

}

type listDockerContainerOutput struct {
	ID    string `json:"ID"`
	Image string `json:"Image"`
	State string `json:"State"`
	Names string `json:"Names"`
}
