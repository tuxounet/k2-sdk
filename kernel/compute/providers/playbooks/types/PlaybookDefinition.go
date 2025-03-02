package types

type PlaybookDefinition struct {
	Name      string `json:"name"`
	Order     int    `json:"order"`
	Provision string `json:"provision"`
	Start     string `json:"start"`
	Stop      string `json:"stop"`
	Teardown  string `json:"teardown"`
}
