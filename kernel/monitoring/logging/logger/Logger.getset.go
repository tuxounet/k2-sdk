package logger

import "github.com/sirupsen/logrus"

func (l *Logger) GetName() string {
	return l.name
}
func (l *Logger) SetLevel(level logrus.Level) {
	l.log.Logger.SetLevel(level)
}

func (l *Logger) GetBaseLogger() *logrus.Entry {
	return l.log
}
