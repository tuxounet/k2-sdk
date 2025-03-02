package types

import "github.com/sirupsen/logrus"

type ILogger interface {
	GetName() string
	CreateSubLogger(name string) ILogger
	Scope(name string, handler func(log ILogger) error) error
	ScopeWithReturn(name string, handler func(log ILogger) (any, error)) (any, error)
	GetBaseLogger() *logrus.Entry
	Trace(message string)
	TraceF(format string, args ...interface{})
	Debug(message string)
	DebugF(format string, args ...interface{})
	Info(message string)
	InfoF(format string, args ...interface{})
	Warn(message string)
	WarnF(format string, args ...interface{})
	Error(message string)
	ErrorF(format string, args ...interface{})
	Panic(message string)
	PanicF(format string, args ...interface{})
}

type ILoggable interface {
	GetLogger() ILogger
}
