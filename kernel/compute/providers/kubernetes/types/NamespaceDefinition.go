package types

import "embed"

type NamespaceDefinition struct {
	Name      string    `json:"name"`
	Order     int       `json:"order"`
	Templates *embed.FS `json:"templates"`
}
