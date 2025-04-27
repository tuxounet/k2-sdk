package types

import (
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type IngressDefinition struct {
	IngressPath  string                     `json:"path"`
	IngressHost  string                     `json:"host"`
	AccessPolicy runtimeTypes.IAccessPolicy `json:"accessPolicy"`
	ServicePort  int                        `json:"servicePort"`
	ServiceHost  string                     `json:"serviceHost"`
}

type IngressRegistarFunction func(ingress *IngressDefinition) error
