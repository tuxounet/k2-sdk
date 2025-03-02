package containers

import (
	"github.com/tuxounet/k2-sdk/system"
)

func (p *Provider) withPodmanCliCmd(args ...string) *system.CmdCall {

	log := p.GetLogger().CreateSubLogger("exec")
	cmdCall := system.NewCmdCall(log, "podman", args...)

	runDir := p.GetService().GetKernel().GetRunDirectory()
	cmdCall.Cwd = &runDir

	return cmdCall

}

func (p *Provider) listContainers() ([]*listContainersContainerOutput, error) {

	out, err := system.JsonCommandOutput[[]map[string]interface{}](p.withPodmanCliCmd("ps", "--all", "--format", "json"))
	if err != nil {
		p.GetLogger().ErrorF("%s provider listContainers failed: %s", ProviderKey, err)
		return nil, err
	}

	ret := []*listContainersContainerOutput{}

	for _, container := range out {
		ret = append(ret, &listContainersContainerOutput{
			ID:    container["Id"].(string),
			Image: container["Image"].(string),
			State: container["State"].(string),
			Names: container["Names"].([]interface{}),
		})

	}

	return ret, nil

}

type listContainersContainerOutput struct {
	ID    string        `json:"ID"`
	Image string        `json:"Image"`
	State string        `json:"State"`
	Names []interface{} `json:"Names"`
}

func (p *Provider) killRootlessPort() {
	processName := "rootlessport"
	p.GetLogger().TraceF("trying killing process: %s", processName)
	killCmd := system.NewCmdCall(p.GetLogger(), "pkill", "-f", processName)

	err := system.OsExec(killCmd)
	if err != nil {
		p.GetLogger().WarnF("Failed to kill process %s: %s", processName, err)
	}

}
