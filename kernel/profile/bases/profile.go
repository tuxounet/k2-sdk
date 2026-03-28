package bases

type ProfilePublic struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Profile struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	Properties map[string]string `json:"properties"`
	Secrets    map[string]string `json:"secrets"`
}

func (p *Profile) Public() *ProfilePublic {
	return &ProfilePublic{
		Name:    p.Name,
		Version: p.Version,
	}
}
