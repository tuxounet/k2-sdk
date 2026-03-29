package logger

import (
	"context"
	"os"

	"github.com/tuxounet/k2-sdk/kernel/monitoring/logging/logger/hooks"
	"github.com/tuxounet/k2-sdk/types"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	name string
	log  *logrus.Entry
}

func NewRootLogger(ctx context.Context, name string, logFolder string, level logrus.Level) *Logger {
	rootLogger := logrus.New()

	// Auto-detect systemd/journald: if JOURNAL_STREAM is set, timestamps are redundant
	_, underJournald := os.LookupEnv("JOURNAL_STREAM")

	rootLogger.SetFormatter(&ConsoleFormatter{
		DisableTimestamp: underJournald,
		TimestampFormat:  "15:04:05.000",
		AppName:          name,
	})
	rootLogger.SetLevel(level)

	err := os.MkdirAll(logFolder, 0755)
	if err != nil {
		rootLogger.Warnf("Failed to create log folder %s", logFolder)
	}

	fileWriterHook := hooks.NewFileWriterHook(logFolder, name)
	rootLogger.Hooks.Add(fileWriterHook)

	log := rootLogger.WithContext(ctx).WithField("log", name)

	return &Logger{
		name: name,
		log:  log,
	}
}

func (l *Logger) CreateSubLogger(name string) types.ILogger {

	loggerName := l.name + "." + name
	log := l.log.WithField("log", loggerName)

	return &Logger{
		name: loggerName,
		log:  log,
	}
}

func (l *Logger) Scope(name string, handler func(log types.ILogger) error) error {

	log := l.CreateSubLogger(name)
	log.TraceF("BEGIN")

	err := handler(log)
	if err != nil {

		log.ErrorF("FAULTED: %s", name, err.Error())
		return err
	}

	log.TraceF("END")

	return nil

}

func (l *Logger) ScopeWithReturn(name string, handler func(log types.ILogger) (any, error)) (any, error) {

	log := l.CreateSubLogger(name)
	log.TraceF("BEGIN")

	ret, err := handler(log)
	if err != nil {

		log.ErrorF("FAULTED: %s", name, err.Error())
		return nil, err
	}

	log.TraceF("END")

	return ret, nil

}
