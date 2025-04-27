package types

import (
	ingressTypes "github.com/tuxounet/k2-sdk/kernel/network/ingress/types"
)

type PortsMapRecord struct {
	LocalPort     int                             `json:"localPort"`
	Order         int                             `json:"index"`
	ContainerName string                          `json:"name"`
	ContainerPort int                             `json:"containerPort"`
	Ingress       *ingressTypes.IngressDefinition `json:"ingress"`
}
