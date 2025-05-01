package types

import (
	"github.com/gin-gonic/gin"
	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

type IngressDefinition struct {
	IngressPath   string                     `json:"path"`
	AccessPolicy  runtimeTypes.IAccessPolicy `json:"accessPolicy"`
	ServicePort   int                        `json:"servicePort"`
	ServiceHost   string                     `json:"serviceHost"`
	CustomHandler gin.HandlerFunc            `json:"-"`
}

type IngressRegistarFunction func(ingress *IngressDefinition) error
