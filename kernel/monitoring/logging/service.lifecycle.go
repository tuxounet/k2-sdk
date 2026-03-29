package logging

import (
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tuxounet/k2-sdk/kernel/monitoring/logging/logger"
)

func parseLogLevel(levelStr string) logrus.Level {
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func (s *LoggingService) Init() error {

	if s.GetData("rootLogger") == nil {
		kernel := s.GetKernel()
		logsFolder := path.Join(kernel.GetRunDirectory(), "var", "log")

		// Determine log level: env var K2_HOST_LOGGING_LEVEL, default to debug
		level := logrus.DebugLevel
		if envLevel := os.Getenv("K2_HOST_LOGGING_LEVEL"); envLevel != "" {
			level = parseLogLevel(envLevel)
		}

		rootLogger := logger.NewRootLogger(kernel.GetRootContext(), kernel.GetApp().GetName(), logsFolder, level)
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
