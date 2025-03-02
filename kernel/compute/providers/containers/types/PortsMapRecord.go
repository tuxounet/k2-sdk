package types

type PortsMapRecord struct {
	LocalPort     int                         `json:"localPort"`
	Order         int                         `json:"index"`
	ContainerName string                      `json:"name"`
	Ingress       *ContainerDefinitionIngress `json:"ingress"`
}
