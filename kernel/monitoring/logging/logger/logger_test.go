package logger

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuxounet/k2-sdk/types"
)

func newTestLogger(t *testing.T) (*Logger, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "k2-logger-test-*")
	require.NoError(t, err)

	logger := NewRootLogger(context.Background(), "test", tmpDir, logrus.TraceLevel)
	return logger, tmpDir
}

func TestNewRootLogger(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	assert.NotNil(t, logger)
	assert.Equal(t, "test", logger.GetName())
	assert.NotNil(t, logger.GetBaseLogger())
}

func TestLogger_CreateSubLogger(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	sub := logger.CreateSubLogger("child")
	assert.NotNil(t, sub)

	subLogger := sub.(*Logger)
	assert.Equal(t, "test.child", subLogger.GetName())
}

func TestLogger_NestedSubLoggers(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	sub1 := logger.CreateSubLogger("a")
	sub2 := sub1.CreateSubLogger("b")

	assert.Equal(t, "test.a.b", sub2.(*Logger).GetName())
}

func TestLogger_AllMethods(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	// These should not panic
	logger.Trace("trace message")
	logger.TraceF("trace %s", "format")
	logger.Debug("debug message")
	logger.DebugF("debug %s", "format")
	logger.Info("info message")
	logger.InfoF("info %s", "format")
	logger.Warn("warn message")
	logger.WarnF("warn %s", "format")
	logger.Error("error message")
	logger.ErrorF("error %s", "format")
}

func TestLogger_SetLevel(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	logger.SetLevel(logrus.ErrorLevel)
	assert.Equal(t, logrus.ErrorLevel, logger.log.Logger.GetLevel())
}

func TestLogger_Scope_Success(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	executed := false
	err := logger.Scope("scoped", func(log types.ILogger) error {
		executed = true
		assert.Equal(t, "test.scoped", log.GetName())
		return nil
	})
	require.NoError(t, err)
	assert.True(t, executed)
}

func TestLogger_Scope_Error(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	err := logger.Scope("scoped", func(log types.ILogger) error {
		return errors.New("test error")
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test error")
	assert.Contains(t, err.Error(), "scope scoped")
	assert.ErrorIs(t, err, errors.Unwrap(err))
}

func TestLogger_ScopeWithReturn_Success(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	result, err := logger.ScopeWithReturn("scoped", func(log types.ILogger) (any, error) {
		return "result-value", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "result-value", result)
}

func TestLogger_ScopeWithReturn_Error(t *testing.T) {
	logger, tmpDir := newTestLogger(t)
	defer os.RemoveAll(tmpDir)

	result, err := logger.ScopeWithReturn("scoped", func(log types.ILogger) (any, error) {
		return nil, errors.New("scope error")
	})
	assert.Error(t, err)
	assert.Nil(t, result)
}
