package ansible

import (
	playbooksBases "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/bases"
	playbooksTypes "github.com/tuxounet/k2-sdk/kernel/compute/providers/playbooks/types"

	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	playbooksBases.BaseControllerPlaybook
}

const ControllerKey = "ansible"

func NewController(component types.IAppComponent) types.IAppController {

	base := playbooksBases.NewBaseControllerPlaybook(component, &playbooksTypes.PlaybookDefinition{
		Name:  ControllerKey,
		Order: 3,
		Start: `
- name: Hello world
  debug:
    msg: "hello: {{ messages.french.hello }}"
`,
		Stop: `
- name: Goodbye world
  debug:
    msg: "goodbye: {{ messages.french.goodbye }}"
`,
	}, nil)
	return &Controller{
		base,
	}
}
