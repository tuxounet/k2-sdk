package types

import (
	"embed"

	"github.com/tuxounet/k2-sdk/types"
)

type NamespaceDefinition struct {
	Name      string             `json:"name"`
	Order     int                `json:"order"`
	Templates *embed.FS          `json:"templates"`
	Ingresses []NamespaceIngress `json:"ingresses"`
}

type NamespaceIngress struct {
	AccessPolicy types.IAccessPolicy `json:"accessPolicy"`
	IngressPath  string              `json:"ingressPath"`
	ServiceName  string              `json:"serviceName"`
	ServicePort  int                 `json:"servicePort"`
}
