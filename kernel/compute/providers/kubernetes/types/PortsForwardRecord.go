package types

type PortsForwardRecord struct {
	Path             string `json:"path"`
	LocalPort        int    `json:"localPort"`
	ServiceNamespace string `json:"serviceNamespace"`
	ServiceName      string `json:"serviceName"`
	ServicePort      int    `json:"servicePort"`
}
