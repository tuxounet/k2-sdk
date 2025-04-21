package bases

import (
	"github.com/tuxounet/k2-sdk/kernel/compute/types"
)

func (p *BasePlateformProvider[D]) Init() error {
	return nil
}

func (p *BasePlateformProvider[D]) Render() ([]types.RunnerDefinition, error) {
	return nil, nil
}
func (p *BasePlateformProvider[D]) Setup() error {
	return nil
}
