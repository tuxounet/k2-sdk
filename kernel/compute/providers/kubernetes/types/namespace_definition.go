package types

import (
	"embed"

	"github.com/tuxounet/k2-sdk/types"
)

type NamespaceDefinition struct {
	Name        string                  `json:"name"`
	Order       int                     `json:"order"`
	Templates   *embed.FS               `json:"templates"`
	Ingresses   []NamespaceIngress      `json:"ingresses"`
	Ports       []NamespacePortForwards `json:"ports"`
	WaitTimeout int                     `json:"waitTimeout"`
}

type NamespaceIngress struct {
	AccessPolicy     types.IAccessPolicy `json:"accessPolicy"`
	IngressPath      string              `json:"ingressPath"`
	RewritePath      *string             `json:"rewritePath,omitempty"`
	ServiceNamespace string              `json:"serviceNamespace"`
	ServiceName      string              `json:"serviceName"`
	ServicePort      int                 `json:"servicePort"`
}

type NamespacePortForwards struct {
	LocalPort        int    `json:"localPort"`
	ServiceNamespace string `json:"serviceNamespace"`
	ServiceName      string `json:"serviceName"`
	ServicePort      int    `json:"servicePort"`
}
