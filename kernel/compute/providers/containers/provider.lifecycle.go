package containers

import (
	"fmt"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/engines"
	containersTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/containers/types"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
)

func (p *Provider) Init() error {
	err := p.getPortsMapSore().Nuke()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke portsmap store: %s", err)
		return err
	}

	return nil
}

func (p *Provider) Nuke() error {
	p.GetLogger().TraceF("Nuking up %s provider", ProviderKey)

	containerEngine := p.getContainerEngine()
	if containerEngine == nil {
		return fmt.Errorf("container engine is not set")
	}
	err := containerEngine.Nuke()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke %s provider: %s", ProviderKey, err)
		return err
	}

	return nil
}

func (p *Provider) Setup() error {
	p.GetLogger().TraceF("Setting up %s provider", ProviderKey)

	containerEngine := p.getContainerEngine()
	if containerEngine == nil {
		return fmt.Errorf("container engine is not set")
	}
	engine := containerEngine
	err := engine.Setup()
	if err != nil {
		p.GetLogger().ErrorF("Failed to nuke %s provider: %s", ProviderKey, err)
		return err
	}

	p.GetLogger().DebugF("%s provider setup done", ProviderKey)
	return nil

}

func (p *Provider) Render() ([]computeTypes.RunnerDefinition, error) {

	p.GetLogger().DebugF("[RENDER] Rendering %s provider", ProviderKey)

	definitions := p.GetDefinitions()
	engineName := p.getEngine()
	var engine containersTypes.IContainerEngine
	switch engineName {
	case "podman":
		engine = engines.NewPodmanEngine(p.GetService(), p.getPortsMapSore(), p.GetIngressRegistar())
	case "docker":
		engine = engines.NewDockerEngine(p.GetService(), p.getPortsMapSore(), p.GetIngressRegistar())
	default:
		return nil, fmt.Errorf("unknown engine: %s", engine)
	}
	p.setContainerEngine(engine)

	runners := make([]computeTypes.RunnerDefinition, 0)

	for _, definition := range definitions {

		newRunnerDefinition := computeTypes.RunnerDefinition{
			Name:      definition.Name,
			Order:     definition.Order,
			Provider:  ProviderKey,
			Provision: "",
			Start:     "",
			Stop:      "",
			Teardown:  "",
		}

		provisionScript, err := engine.RenderPlaybookTasks(definition, "provision")
		if err != nil {
			return nil, err
		}
		newRunnerDefinition.Provision = provisionScript

		startScript, err := engine.RenderPlaybookTasks(definition, "start")
		if err != nil {
			return nil, err
		}
		newRunnerDefinition.Start = startScript

		stopScript, err := engine.RenderPlaybookTasks(definition, "stop")
		if err != nil {
			return nil, err
		}
		newRunnerDefinition.Stop = stopScript

		teardownScript, err := engine.RenderPlaybookTasks(definition, "teardown")
		if err != nil {
			return nil, err
		}
		newRunnerDefinition.Teardown = teardownScript

		runners = append(runners, newRunnerDefinition)
	}

	p.GetLogger().DebugF("Rendered %d runners for %s provider", len(runners), ProviderKey)

	return runners, nil
}
