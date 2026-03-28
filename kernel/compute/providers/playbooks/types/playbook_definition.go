package types

type PlaybookFileDefinition struct {
	Content string `json:"content"`
	Mode    string `json:"mode"`
}

type PlaybookDefinition struct {
	Name      string                            `json:"name"`
	Order     int                               `json:"order"`
	Provision string                            `json:"provision"`
	Start     string                            `json:"start"`
	Stop      string                            `json:"stop"`
	Teardown  string                            `json:"teardown"`
	Files     map[string]PlaybookFileDefinition `json:"files"`
}
