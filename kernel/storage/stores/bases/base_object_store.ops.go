package bases

import "github.com/tuxounet/k2-sdk/system"

func (s *BaseObjectStore[R]) HasValue() (bool, error) {
	store, err := s.getStore()
	if err != nil {
		s.log.ErrorF("Failed to read state: %s", err)
		return false, err
	}

	exist, err := store.Exists(s.valueKey)
	if err != nil {
		s.log.ErrorF("Failed to check state: %s", err)
		return false, err
	}
	return exist, nil

}

func (s *BaseObjectStore[R]) GetValue() (*R, error) {

	hasState, err := s.HasValue()
	if err != nil {
		s.log.ErrorF("Failed to check state: %s", err)
		return nil, err
	}

	currentValue := s.defaultValue

	if hasState {

		store, err := s.getStore()
		if err != nil {
			s.log.ErrorF("Failed to read state: %s", err)
			return nil, err
		}

		value, err := store.ReadObject(s.valueKey)
		if err != nil {
			s.log.ErrorF("Failed to read state: %s", err)
			return nil, err
		}
		currentValue = string(value)

	}

	result, err := system.LoadJSONFromString[R](currentValue)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *BaseObjectStore[R]) SetValue(newValue R) error {

	store, err := s.getStore()
	if err != nil {
		s.log.ErrorF("Failed to set value: %s", err)
		return err
	}

	rawValue, err := system.DumpToJsonString(newValue)
	if err != nil {
		s.log.ErrorF("Failed to unmarshall value")
		return err
	}

	err = store.WriteObject(s.valueKey, []byte(rawValue))
	if err != nil {
		s.log.ErrorF("Failed to write value: %s", err)
		return err
	}

	return nil
}

func (s *BaseObjectStore[R]) Nuke() error {

	exists, err := s.HasValue()

	if err != nil {
		s.GetLogger().ErrorF("Failed to check state: %s", err)
		return err
	}
	if !exists {
		return nil
	}

	store, err := s.getStore()
	if err != nil {
		s.log.ErrorF("Failed to read state: %s", err)
		return err
	}

	err = store.DeleteObject(s.valueKey)
	if err != nil {
		s.log.ErrorF("Failed to delete state: %s", err)
		return err
	}

	return nil

}
