package logging

import (
	"github.com/tuxounet/k2-sdk/types"
)

const LoggingServiceKey = "monitoring.logging"

type LoggingService struct {
	name   string
	kernel types.IKernel

	config map[string]string
	data   map[string]interface{}
}

func NewService(k types.IKernel) *LoggingService {

	instance := &LoggingService{
		name:   LoggingServiceKey,
		kernel: k,
		config: make(map[string]string),
		data:   make(map[string]interface{}),
	}

	return instance

}

func (s *LoggingService) GetRootLogger() types.ILogger {
	return s.GetData("rootLogger").(types.ILogger)
}

func (s *LoggingService) GetLogger() types.ILogger {
	return s.GetData("serviceLogger").(types.ILogger)
}
