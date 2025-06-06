package logger

func (l *Logger) Trace(message string) {
	l.log.Tracef("➡️  %s > %s", l.name, message)
}

func (l *Logger) TraceF(format string, args ...interface{}) {
	args = append([]interface{}{l.name}, args...)
	l.log.Tracef("➡️  %s > "+format, args...)
}

func (l *Logger) Debug(message string) {
	l.log.Debugf("▶️  %s > %s", l.name, message)
}

func (l *Logger) DebugF(format string, args ...interface{}) {

	args = append([]interface{}{l.name}, args...)
	l.log.Debugf("▶️  %s > "+format, args...)
}

func (l *Logger) Info(message string) {
	l.log.Infof("🌟 %s > %s", l.name, message)
}

func (l *Logger) InfoF(format string, args ...interface{}) {
	args = append([]interface{}{l.name}, args...)
	l.log.Infof("🌟 %s > "+format, args...)

}

func (l *Logger) Warn(message string) {
	l.log.Warnf("🫵  %s > %s", l.name, message)
}

func (l *Logger) WarnF(format string, args ...interface{}) {
	args = append([]interface{}{l.name}, args...)
	l.log.Warnf("🫵  %s > "+format, args...)

}

func (l *Logger) Error(message string) {
	l.log.Errorf("🛑 %s > %s", l.name, message)
}

func (l *Logger) ErrorF(format string, args ...interface{}) {
	args = append([]interface{}{l.name}, args...)
	l.log.Errorf("🛑 %s > "+format, args...)
}

func (l *Logger) Panic(message string) {
	l.log.Panicf("🔥 %s > %s", l.name, message)
}

func (l *Logger) PanicF(format string, args ...interface{}) {
	args = append([]interface{}{l.name}, args...)
	l.log.Panicf("🔥 %s > "+format, args...)
}
