package logging

import (
	"path"

	"github.com/sirupsen/logrus"
	"github.com/tuxounet/k2-sdk/kernel/monitoring/logging/logger"
)

func (s *LoggingService) Init() error {

	if s.GetData("rootLogger") == nil {
		kernel := s.GetKernel()
		logsFolder := path.Join(kernel.GetRunDirectory(), "var", "log")
		rootLogger := logger.NewRootLogger(kernel.GetRootContext(), kernel.GetApp().GetName(), logsFolder, logrus.TraceLevel)
		s.SetData("rootLogger", rootLogger)
		s.GetRootLogger().Debug("Root logger initialized")
	}
	if s.GetData("serviceLogger") == nil {
		serviceLogger := s.GetRootLogger().CreateSubLogger(s.GetName())
		s.SetData("serviceLogger", serviceLogger)

		s.GetLogger().Debug("Service logger initialized")
	}
	return nil
}

func (s *LoggingService) Register() error {
	return nil
}

func (s *LoggingService) Start() error {
	return nil
}

func (s *LoggingService) Stop() error {
	return nil
}
