package types

type Store struct {
	Name            string                 `json:"name"`
	Backend         string                 `json:"backend"`
	Flags           map[string]interface{} `json:"flags"`
	backendProvider IStoreProvider
}

func NewStore(name string, backend string, flags map[string]interface{}) *Store {
	return &Store{
		Name:    name,
		Backend: backend,
		Flags:   flags,
	}
}

func (s *Store) ResolveBackend(backends []IStoreProvider) error {
	for _, backend := range backends {
		if backend.GetName() == s.Backend {
			s.backendProvider = backend
			return nil
		}
	}
	return nil

}

func (s *Store) GetName() string {
	return s.Name
}

func (s *Store) Exists(key string) (bool, error) {
	return s.backendProvider.Exists(s, key)
}

func (s *Store) ReadObject(key string) ([]byte, error) {
	return s.backendProvider.Read(s, key)
}

func (s *Store) WriteObject(key string, data []byte) error {
	return s.backendProvider.Write(s, key, data)
}

func (s *Store) DeleteObject(key string) error {
	return s.backendProvider.Delete(s, key)
}
