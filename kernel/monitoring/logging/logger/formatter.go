package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var levelEmojis = map[logrus.Level]string{
	logrus.TraceLevel: "➡️ ",
	logrus.DebugLevel: "▶️ ",
	logrus.InfoLevel:  "🌟",
	logrus.WarnLevel:  "🫵 ",
	logrus.ErrorLevel: "🛑",
	logrus.PanicLevel: "🔥",
	logrus.FatalLevel: "🔥",
}

var levelLabels = map[logrus.Level]string{
	logrus.TraceLevel: "TRC",
	logrus.DebugLevel: "DBG",
	logrus.InfoLevel:  "INF",
	logrus.WarnLevel:  "WRN",
	logrus.ErrorLevel: "ERR",
	logrus.PanicLevel: "PNC",
	logrus.FatalLevel: "FTL",
}

// ConsoleFormatter formats log entries for human-readable console output.
// It uses emojis, short scope names, and indentation based on scope depth.
type ConsoleFormatter struct {
	// DisableTimestamp removes the timestamp from output (useful under systemd/journald).
	DisableTimestamp bool
	// TimestampFormat is the time format string (Go reference time layout).
	TimestampFormat string
	// AppName is the root application name, stripped from scope display for brevity.
	AppName string
}

func (f *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	// Timestamp
	if !f.DisableTimestamp {
		ts := entry.Time.Format(f.timestampFormat())
		b.WriteString(ts)
		b.WriteByte(' ')
	}

	// Level label + emoji
	label := levelLabels[entry.Level]
	emoji := levelEmojis[entry.Level]
	fmt.Fprintf(&b, "[%s] %s ", label, emoji)

	// Scope from the "log" field
	scope := f.extractScope(entry)
	if scope != "" {
		depth := strings.Count(scope, ".")
		if depth > 0 {
			b.WriteString(strings.Repeat("  ", depth))
		}
		fmt.Fprintf(&b, "%s > ", scope)
	}

	// Message
	b.WriteString(entry.Message)
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *ConsoleFormatter) timestampFormat() string {
	if f.TimestampFormat != "" {
		return f.TimestampFormat
	}
	return "15:04:05.000"
}

func (f *ConsoleFormatter) extractScope(entry *logrus.Entry) string {
	logField, ok := entry.Data["log"]
	if !ok {
		return ""
	}
	scope, ok := logField.(string)
	if !ok {
		return ""
	}
	// Strip the app name prefix for brevity
	if f.AppName != "" && strings.HasPrefix(scope, f.AppName) {
		scope = strings.TrimPrefix(scope, f.AppName)
		scope = strings.TrimPrefix(scope, ".")
	}
	return scope
}

// FileFormatter formats log entries for file output: no emojis, no colors,
// always includes full timestamp and level text.
type FileFormatter struct {
	TimestampFormat string
}

func (f *FileFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	ts := entry.Time.Format(f.timestampFormat())
	fmt.Fprintf(&b, "%s ", ts)

	label := levelLabels[entry.Level]
	fmt.Fprintf(&b, "[%s] ", label)

	// Scope
	if logField, ok := entry.Data["log"]; ok {
		if scope, ok := logField.(string); ok {
			fmt.Fprintf(&b, "%s > ", scope)
		}
	}

	b.WriteString(entry.Message)
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *FileFormatter) timestampFormat() string {
	if f.TimestampFormat != "" {
		return f.TimestampFormat
	}
	return "2006-01-02 15:04:05.000"
}
