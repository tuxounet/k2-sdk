package types

type RunnerVerb string

const (
	RunnerVerbProvision RunnerVerb = "provision"
	RunnerVerbStart     RunnerVerb = "start"
	RunnerVerbStop      RunnerVerb = "stop"
	RunnerVerbTeardown  RunnerVerb = "teardown"
)

type RunnerDefinition struct {
	Order     int    `json:"order"`
	Provider  string `json:"plateform"`
	Name      string `json:"name"`
	Provision string `json:"provision"`
	Teardown  string `json:"teardown"`
	Start     string `json:"start"`
	Stop      string `json:"stop"`
}
