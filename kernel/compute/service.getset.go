package compute

import (
	"sort"

	"github.com/tuxounet/k2-sdk/kernel/compute/types"
)

func (s *Service) getProviders() []types.IBasePlateformProvider {
	return s.GetData("providers").([]types.IBasePlateformProvider)
}

func (s *Service) GetProvider(name string) types.IBasePlateformProvider {
	providers := s.getProviders()
	for _, provider := range providers {
		if provider.GetName() == name {
			return provider
		}

	}
	return nil
}

func (s *Service) setProviders(providers []types.IBasePlateformProvider) {
	s.SetData("providers", providers)
}

func (s *Service) getRunners() []types.RunnerDefinition {
	allRunners := s.GetData("runners").([]types.RunnerDefinition)

	orderedRunners := make([]types.RunnerDefinition, len(allRunners))
	copy(orderedRunners, allRunners)
	sort.SliceStable(orderedRunners, func(i, j int) bool {
		return orderedRunners[i].Order < orderedRunners[j].Order
	})

	return orderedRunners
}
func (s *Service) getReverseRunners() []types.RunnerDefinition {
	runners := s.getRunners()
	reverseRunners := make([]types.RunnerDefinition, len(runners))
	for i, j := 0, len(runners)-1; i < len(runners); i, j = i+1, j-1 {
		reverseRunners[i] = runners[j]
	}
	return reverseRunners
}

func (s *Service) setRunners(runners []types.RunnerDefinition) {
	s.SetData("runners", runners)
}
