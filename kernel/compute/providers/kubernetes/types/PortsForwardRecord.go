package types

type PortsForwardRecord struct {
	LocalPort        int    `json:"localPort"`
	ServiceNamespace string `json:"serviceNamespace"`
	ServiceName      string `json:"serviceName"`
	ServicePort      int    `json:"servicePort"`
}
