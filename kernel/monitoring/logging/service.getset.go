package logging

import (
	"github.com/tuxounet/k2-sdk/types"
)

func (s *LoggingService) GetName() string {
	return s.name
}

func (s *LoggingService) GetKernel() types.IKernel {
	return s.kernel
}
func (s *LoggingService) GetConfig(key string) string {
	return s.config[key]
}

func (s *LoggingService) SetConfig(key string, value string) {
	s.config[key] = value
}

func (s *LoggingService) GetData(key string) interface{} {
	return s.data[key]
}

func (s *LoggingService) SetData(key string, value interface{}) {
	s.data[key] = value
}
