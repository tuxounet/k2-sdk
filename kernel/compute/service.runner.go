package compute

import (
	"fmt"
	"strings"

	"github.com/tuxounet/k2-sdk/kernel/compute/types"
	"github.com/tuxounet/k2-sdk/system"
)

func (s *Service) resetRunners() error {
	allRunners := make([]types.RunnerDefinition, 0)
	s.setRunners(allRunners)
	return nil
}

func (s *Service) renderProvisionRunners() error {

	verb := types.RunnerVerbProvision
	tasks := ""
	allRunners := s.getRunners()

	for _, r := range allRunners {
		if r.Provision != "" {
			tasks += r.Provision + "\n"
		}
	}

	err := s.renderPlaybook(verb, tasks)
	if err != nil {
		return fmt.Errorf("failed to render %s playbook : %s", verb, err.Error())
	}

	return nil

}

func (s *Service) renderStartRunners() error {
	verb := types.RunnerVerbStart
	tasks := ""
	allRunners := s.getRunners()

	for _, r := range allRunners {
		if r.Start != "" {
			tasks += r.Start + "\n"
		}
	}

	err := s.renderPlaybook(verb, tasks)
	if err != nil {
		return fmt.Errorf("failed to render %s playbook : %s", verb, err.Error())
	}

	return nil
}

func (s *Service) renderStopRunners() error {
	verb := types.RunnerVerbStop
	tasks := ""
	allRunners := s.getReverseRunners()

	for _, r := range allRunners {
		if r.Stop != "" {
			tasks += r.Stop + "\n"
		}
	}

	err := s.renderPlaybook(verb, tasks)
	if err != nil {
		return fmt.Errorf("failed to render %s playbook : %s", verb, err.Error())
	}

	return nil

}

func (s *Service) renderTeardownRunners() error {
	verb := types.RunnerVerbTeardown
	tasks := ""
	allRunners := s.getReverseRunners()

	for _, r := range allRunners {
		if r.Teardown != "" {
			tasks += r.Teardown + "\n"
		}
	}

	err := s.renderPlaybook(verb, tasks)
	if err != nil {
		return fmt.Errorf("failed to render %s playbook : %s", verb, err.Error())
	}
	return nil
}

func (s *Service) renderPlaybook(verb types.RunnerVerb, tasks string) error {

	playbook := ""
	playbook += fmt.Sprintf("- name: %s\n", verb)
	playbook += fmt.Sprintf("  hosts: %s\n", "localhost")
	playbook += fmt.Sprintf("  gather_facts: %s\n", "no")

	if tasks == "" {
		playbook += "  tasks: []\n"
	} else {
		playbook += "  tasks:\n"
		lines := strings.Split(tasks, "\n")
		for _, line := range lines {
			playbook += fmt.Sprintf("    %s\n", line)
		}
	}

	paths := s.getPathsService()
	targetPlaybookFileName := paths.CominePath("etc", "compute", fmt.Sprintf("%s.yaml", verb))
	rootStore, err := s.getRootStore()
	if err != nil {
		return fmt.Errorf("failed to retrieve root store : %s", err.Error())

	}
	err = rootStore.WriteObject(targetPlaybookFileName, []byte(playbook))
	if err != nil {
		return fmt.Errorf("failed to wrte %s playbook : %s", verb, err.Error())

	}
	return nil

}

func (s *Service) execPlaybook(verb types.RunnerVerb) error {

	paths := s.getPathsService()
	playbookPath := paths.CominePath("etc", "compute", fmt.Sprintf("%s.yaml", verb))
	inventoryPath := paths.CominePath("etc", "compute", "inventory")
	cwd := s.GetKernel().GetRunDirectory()
	cmdCall := system.NewCmdCall(s.GetLogger(), "ansible-playbook", "-i", inventoryPath, playbookPath, "--extra-vars", fmt.Sprintf("run_dir=%s", s.GetKernel().GetRunDirectory()))
	cmdCall.Cwd = &cwd

	_, err := system.OsExecAndTailToLog(cmdCall)
	if err != nil {
		return fmt.Errorf("failed to exec %s: %s", verb, err.Error())
	}

	return nil
}
