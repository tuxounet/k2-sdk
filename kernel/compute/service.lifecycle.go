package compute

import (
	"fmt"

	"github.com/tuxounet/k2-sdk/kernel/compute/providers/containers"
	"github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks"

	"github.com/tuxounet/k2-sdk/kernel/compute/types"
)

func (s *Service) Register() error {

	enabled, err := s.isEnabled()
	if err != nil {
		return fmt.Errorf("unable to check compute service enabled status: %s", err.Error())
	}
	if !enabled {
		s.GetLogger().DebugF("compute service is disabled")
		return nil
	}
	s.GetLogger().TraceF("begin register")
	providers := []types.IBasePlateformProvider{
		containers.NewProvider(s),
		playbooks.NewProvider(s),
	}

	s.setProviders(providers)

	for _, p := range providers {
		err := p.Nuke()
		if err != nil {
			s.GetLogger().ErrorF("provider %s nuke failed: %s", p.GetName(), err)
			return err
		}

		err = p.Setup()
		if err != nil {
			s.GetLogger().ErrorF("provider %s setup failed: %s", p.GetName(), err)
			return err
		}
	}

	err = s.nukeInventory()
	if err != nil {
		s.GetLogger().ErrorF("nuke inventory failed: %s", err)
		return err
	}

	err = s.resetRunners()
	if err != nil {
		return fmt.Errorf("unable ot reset runners collection : %s", err.Error())
	}

	allRunners := make([]types.RunnerDefinition, 0)
	for _, p := range providers {
		runners, err := p.Render()

		if err != nil {
			s.GetLogger().ErrorF("provider %s render failed: %s", p.GetName(), err)
			return err
		}
		allRunners = append(allRunners, runners...)
	}
	s.setRunners(allRunners)

	err = s.renderInventory()
	if err != nil {
		return fmt.Errorf("failed to render inventory: %s ", err.Error())
	}
	err = s.renderProvisionRunners()
	if err != nil {
		return fmt.Errorf("failed to render provision playbook from rendered runners: %s ", err.Error())
	}
	err = s.renderStartRunners()
	if err != nil {
		return fmt.Errorf("failed to render start playbook from rendered runners: %s ", err.Error())
	}
	err = s.renderStopRunners()
	if err != nil {
		return fmt.Errorf("failed to render stop playbook from rendered runners: %s ", err.Error())
	}
	err = s.renderTeardownRunners()
	if err != nil {
		return fmt.Errorf("failed to render teardown playbook from rendered runners: %s ", err.Error())
	}

	s.GetLogger().TraceF("end register")

	return nil
}

func (s *Service) Start() error {
	enabled, err := s.isEnabled()
	if err != nil {
		return fmt.Errorf("unable to check compute service enabled status: %s", err.Error())
	}
	if !enabled {
		s.GetLogger().DebugF("compute service is disabled")
		return nil
	}
	err = s.execPlaybook(types.RunnerVerbProvision)
	if err != nil {
		return fmt.Errorf("provision phase failed: %s", err.Error())
	}

	err = s.execPlaybook(types.RunnerVerbStart)
	if err != nil {
		return fmt.Errorf("start phase failed: %s", err.Error())
	}
	return nil
}

func (s *Service) Stop() error {
	enabled, err := s.isEnabled()
	if err != nil {
		return fmt.Errorf("unable to check compute service enabled status: %s", err.Error())
	}
	if !enabled {
		s.GetLogger().DebugF("compute service is disabled")
		return nil
	}
	s.GetLogger().TraceF("begin stop")
	err = s.execPlaybook(types.RunnerVerbStop)
	if err != nil {
		return fmt.Errorf("stop phase failed: %s", err.Error())
	}

	err = s.execPlaybook(types.RunnerVerbTeardown)
	if err != nil {
		return fmt.Errorf("teardown phase failed: %s", err.Error())
	}

	return nil
}
