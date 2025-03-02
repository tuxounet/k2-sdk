package types

type ContainerDefinition struct {
	Name       string                        `json:"name"`
	Order      int                           `json:"order"`
	Image      string                        `json:"image"`
	Command    *[]string                     `json:"command"`
	Volumes    []*ContainerDefinitionVolume  `json:"volumes"`
	Ingresses  []*ContainerDefinitionIngress `json:"ingresses"`
	Ports      []*ContainerDefinitionPort    `json:"ports"`
	Env        map[string]string             `json:"env"`
	Capacities *[]ContainerCapacity          `json:"capacities"`
	Security   *ContainerDefinitionSecurity  `json:"security"`
}

type ContainerDefinitionSecurity struct {
	Privileged bool   `json:"privileged"`
	RunAsUser  string `json:"runAsUser"`
}

type ContainerDefinitionIngress struct {
	ContainerPort string
	Path          string
}

type ContainerDefinitionPort struct {
	ContainerPort string `json:"containerPort"`
	Protocol      string `json:"protocol"`
	HostAddress   string `json:"hostAddress"`
	HostPort      string `json:"hostPort"`
}

type ContainerDefinitionVolume struct {
	Name          string                           `json:"name,omitempty"`
	ContainerPath string                           `json:"containerPath"`
	Binding       ContainerDefinitionVolumeBinding `json:"target"`
}
type ContainerDefinitionVolumeBindingType string

const (
	ContainerDefinitionVolumeBindingTypeMount   ContainerDefinitionVolumeBindingType = "mount"
	ContainerDefinitionVolumeBindingTypeContent ContainerDefinitionVolumeBindingType = "content"
)

type ContainerDefinitionVolumeBinding struct {
	Type     ContainerDefinitionVolumeBindingType `json:"type"`
	HostPath string                               `json:"hostPath"`
	Content  string                               `json:"content"`
}

type ContainerCapacityAction string

const (
	ContainerCapacityActionAdd    ContainerCapacityAction = "add"
	ContainerCapacityActionRemove ContainerCapacityAction = "remove"
)

type ContainerCapacity struct {
	Name   string                  `json:"name"`
	Action ContainerCapacityAction `json:"action"`
}
