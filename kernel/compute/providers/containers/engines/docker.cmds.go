package engines

import (
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

	out, err := system.LoadJSONFromString[[]map[string]any](raw)
	if err != nil {
		e.logger.ErrorF("%s provider listContainers failed: %s", e.Name, err)
		return nil, err
	}

	ret := []*listDockerContainerOutput{}

	for _, container := range out {
		ret = append(ret, &listDockerContainerOutput{
			ID:    container["Id"].(string),
			Image: container["Image"].(string),
			State: container["State"].(string),
			Names: container["Names"].([]any),
		})

	}

	return ret, nil

}

type listDockerContainerOutput struct {
	ID    string        `json:"ID"`
	Image string        `json:"Image"`
	State string        `json:"State"`
	Names []interface{} `json:"Names"`
}
