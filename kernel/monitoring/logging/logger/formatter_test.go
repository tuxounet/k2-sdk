package logger

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsoleFormatter_Format(t *testing.T) {
	f := &ConsoleFormatter{
		DisableTimestamp: true,
		AppName:          "MyApp",
	}

	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.InfoLevel,
		Message: "hello world",
		Time:    time.Now(),
		Data: logrus.Fields{
			"log": "MyApp.service.sub",
		},
	}

	out, err := f.Format(entry)
	require.NoError(t, err)
	s := string(out)

	// Should contain level label and emoji
	assert.Contains(t, s, "[INF]")
	assert.Contains(t, s, "🌟")
	// Should contain the scope without the app name prefix
	assert.Contains(t, s, "service.sub >")
	// Should contain the message
	assert.Contains(t, s, "hello world")
	// Should NOT contain a timestamp
	assert.NotContains(t, s, "2026")
}

func TestConsoleFormatter_WithTimestamp(t *testing.T) {
	f := &ConsoleFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "15:04:05",
		AppName:          "MyApp",
	}

	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.TraceLevel,
		Message: "trace msg",
		Time:    time.Date(2026, 3, 29, 13, 45, 30, 0, time.UTC),
		Data: logrus.Fields{
			"log": "MyApp",
		},
	}

	out, err := f.Format(entry)
	require.NoError(t, err)
	s := string(out)

	assert.Contains(t, s, "13:45:30")
	assert.Contains(t, s, "[TRC]")
	assert.Contains(t, s, "➡️")
}

func TestConsoleFormatter_Indentation(t *testing.T) {
	f := &ConsoleFormatter{
		DisableTimestamp: true,
		AppName:          "App",
	}

	// depth 0: "App" → stripped to "" → 0 dots → no indent
	entry0 := &logrus.Entry{
		Logger: logrus.New(), Level: logrus.InfoLevel, Message: "root", Time: time.Now(),
		Data: logrus.Fields{"log": "App"},
	}
	out0, _ := f.Format(entry0)

	// depth 1: "App.svc" → "svc" → 0 dots → no indent
	entry1 := &logrus.Entry{
		Logger: logrus.New(), Level: logrus.InfoLevel, Message: "svc", Time: time.Now(),
		Data: logrus.Fields{"log": "App.svc"},
	}
	out1, _ := f.Format(entry1)

	// depth 2: "App.svc.sub" → "svc.sub" → 1 dot → 2 spaces indent
	entry2 := &logrus.Entry{
		Logger: logrus.New(), Level: logrus.InfoLevel, Message: "sub", Time: time.Now(),
		Data: logrus.Fields{"log": "App.svc.sub"},
	}
	out2, _ := f.Format(entry2)

	// Deeper scopes should produce longer lines (more indentation)
	assert.True(t, len(out2) > len(out1))
	assert.True(t, len(out1) >= len(out0))
}

func TestConsoleFormatter_AllLevels(t *testing.T) {
	f := &ConsoleFormatter{DisableTimestamp: true, AppName: "X"}

	levels := []struct {
		level logrus.Level
		label string
	}{
		{logrus.TraceLevel, "TRC"},
		{logrus.DebugLevel, "DBG"},
		{logrus.InfoLevel, "INF"},
		{logrus.WarnLevel, "WRN"},
		{logrus.ErrorLevel, "ERR"},
	}

	for _, tt := range levels {
		entry := &logrus.Entry{
			Logger: logrus.New(), Level: tt.level, Message: "msg", Time: time.Now(),
			Data: logrus.Fields{"log": "X"},
		}
		out, err := f.Format(entry)
		require.NoError(t, err)
		assert.Contains(t, string(out), "["+tt.label+"]")
	}
}

func TestFileFormatter_Format(t *testing.T) {
	f := &FileFormatter{}

	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.WarnLevel,
		Message: "something happened",
		Time:    time.Date(2026, 3, 29, 1, 38, 5, 233000000, time.UTC),
		Data: logrus.Fields{
			"log": "MyApp.network.ingress",
		},
	}

	out, err := f.Format(entry)
	require.NoError(t, err)
	s := string(out)

	assert.Contains(t, s, "2026-03-29 01:38:05.233")
	assert.Contains(t, s, "[WRN]")
	assert.Contains(t, s, "MyApp.network.ingress >")
	assert.Contains(t, s, "something happened")
	// Should NOT contain emojis
	assert.NotContains(t, s, "🫵")
	assert.NotContains(t, s, "➡️")
}
