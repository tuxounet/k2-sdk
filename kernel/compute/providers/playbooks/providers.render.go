package playbooks

import (
	"fmt"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"
	computeTypes "github.com/tuxounet/k2-sdk/kernel/compute/types"
)

func (p *Provider) Render() ([]computeTypes.RunnerDefinition, error) {

	definitions := p.GetDefinitions()
	runners := make([]computeTypes.RunnerDefinition, 0)

	for _, definition := range definitions {
		newRunnerDefinition := computeTypes.RunnerDefinition{
			Order:     definition.Order,
			Name:      definition.Name,
			Plateform: ProviderKey,
			Provision: "",
			Start:     "",
			Stop:      "",
			Teardown:  "",
		}

		if definition.Provision != "" {
			provisionTasks, err := p.renderPlaybookTasks(&definition, computeTypes.RunnerVerbProvision, definition.Provision)
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

		if definition.Teardown != "" {
			teardownTasks, err := p.renderPlaybookTasks(&definition, computeTypes.RunnerVerbTeardown, definition.Teardown)
			if err != nil {
				return nil, err
			}
			newRunnerDefinition.Teardown = teardownTasks

		}
		runners = append(runners, newRunnerDefinition)
	}

	return runners, nil
}

func (p *Provider) renderPlaybookTasks(definition *types.PlaybookDefinition, verb computeTypes.RunnerVerb, script string) (string, error) {
	fullPlaybookTasks := fmt.Sprintf("# [%d] Ansible Playbook %s of %s\n", definition.Order, verb, definition.Name)
	fullPlaybookTasks += script

	return fullPlaybookTasks, nil

}
