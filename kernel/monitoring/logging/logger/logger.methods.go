package logger

import "fmt"

func (l *Logger) Trace(message string) {
	l.log.Trace(message)
}

func (l *Logger) TraceF(format string, args ...interface{}) {
	l.log.Trace(fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(message string) {
	l.log.Debug(message)
}

func (l *Logger) DebugF(format string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Info(message string) {
	l.log.Info(message)
}

func (l *Logger) InfoF(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(message string) {
	l.log.Warn(message)
}

func (l *Logger) WarnF(format string, args ...interface{}) {
	l.log.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Error(message string) {
	l.log.Error(message)
}

func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Panic(message string) {
	l.log.Panic(message)
}

func (l *Logger) PanicF(format string, args ...interface{}) {
	l.log.Panic(fmt.Sprintf(format, args...))
}
