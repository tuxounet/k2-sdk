package types

type PortsMapRecord struct {
	LocalPort     int    `json:"localPort"`
	ContainerName string `json:"name"`
	ContainerPort int    `json:"containerPort"`
}
