package logging

import (
	"context"
	"testing"
	"time"

	"github.com/gatewayd-io/gatewayd/config"
	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

// TestNewHcLogAdapter tests the HcLogAdapter.
func TestNewHcLogAdapter(t *testing.T) {
	consoleOutput := capturer.CaptureStdout(func() {
		logger := NewLogger(
			context.Background(),
			LoggerConfig{
				Output:            []config.LogOutput{config.Console},
				Level:             zerolog.TraceLevel,
				TimeFormat:        zerolog.TimeFormatUnix,
				ConsoleTimeFormat: time.RFC3339,
				NoColor:           true,
			},
		)

		hcLogAdapter := NewHcLogAdapter(&logger, "test")
		hcLogAdapter.SetLevel(hclog.Trace)

		hcLogAdapter.Trace("This is a trace message")
		hcLogAdapter.Debug("This is a debug message")
		hcLogAdapter.Info("This is an info message")
		hcLogAdapter.Warn("This is a warn message")
		hcLogAdapter.Error("This is an error message")
		hcLogAdapter.Log(hclog.Debug, "This is a log message")

		hcLogAdapter.SetLevel(hclog.Warn)
		hcLogAdapter.Trace("This is a trace message, but it should not be logged")
	})

	assert.Contains(t, consoleOutput, "TRC This is a trace message")
	assert.Contains(t, consoleOutput, "DBG This is a debug message")
	assert.Contains(t, consoleOutput, "INF This is an info message")
	assert.Contains(t, consoleOutput, "WRN This is a warn message")
	assert.Contains(t, consoleOutput, "ERR This is an error message")
	assert.Contains(t, consoleOutput, "DBG This is a log message")

	assert.NotContains(t, consoleOutput, "TRC This is a trace message, but it should not be logged")
}

// TestNewHcLogAdapter_LogLevel_Difference tests the HcLogAdapter when the
// logger and the hclog adapter have different log levels.
func TestNewHcLogAdapter_LogLevel_Difference(t *testing.T) {
	consoleOutput := capturer.CaptureStdout(func() {
		logger := NewLogger(
			context.Background(),
			LoggerConfig{
				Output:     []config.LogOutput{config.Console},
				Level:      zerolog.WarnLevel,
				TimeFormat: zerolog.TimeFormatUnix,
				NoColor:    true,
			},
		)

		// The logger is set to WarnLevel, but the hclog adapter is set to TraceLevel.
		// This means that the hclog adapter will log everything, but the logger will
		// only log WarnLevel and above.
		hcLogAdapter := NewHcLogAdapter(&logger, "test")
		hcLogAdapter.SetLevel(hclog.Trace)
		hcLogAdapter.Trace("This is a trace message, but it should not be logged")
		hcLogAdapter.Debug("This is a debug message, but it should not be logged")
		hcLogAdapter.Info("This is an info message, but it should not be logged")
		hcLogAdapter.Warn("This is a warn message")
		hcLogAdapter.Error("This is an error message")
		hcLogAdapter.Log(hclog.Debug, "This is a log message, but it should not be logged")
	})

	assert.NotContains(t, consoleOutput, "TRC This is a trace message, but it should not be logged")
	assert.NotContains(t, consoleOutput, "DBG This is a debug message, but it should not be logged")
	assert.NotContains(t, consoleOutput, "INF This is an info message, but it should not be logged")
	assert.Contains(t, consoleOutput, "WRN This is a warn message")
	assert.Contains(t, consoleOutput, "ERR This is an error message")
	assert.NotContains(t, consoleOutput, "DBG This is a log message, but it should not be logged")
}

// TestNewHcLogAdapter_Log tests the HcLogAdapter.Log method.
func TestNewHcLogAdapter_Log(t *testing.T) {
	consoleOutput := capturer.CaptureStdout(func() {
		logger := NewLogger(
			context.Background(),
			LoggerConfig{
				Output:     []config.LogOutput{config.Console},
				Level:      zerolog.TraceLevel,
				TimeFormat: zerolog.TimeFormatUnix,
				NoColor:    true,
			},
		)

		hcLogAdapter := NewHcLogAdapter(&logger, "test")
		hcLogAdapter.SetLevel(hclog.Trace)

		hcLogAdapter.Log(hclog.Off, "This is a message")
		hcLogAdapter.Log(hclog.NoLevel, "This is yet another message")
		hcLogAdapter.Log(hclog.Trace, "This is a trace message")
		hcLogAdapter.Log(hclog.Debug, "This is a debug message")
		hcLogAdapter.Log(hclog.Info, "This is an info message")
		hcLogAdapter.Log(hclog.Warn, "This is a warn message")
		hcLogAdapter.Log(hclog.Error, "This is an error message")
	})

	assert.NotContains(t, consoleOutput, "This is a message")
	assert.NotContains(t, consoleOutput, "This is yet another message")
	assert.Contains(t, consoleOutput, "TRC This is a trace message")
	assert.Contains(t, consoleOutput, "DBG This is a debug message")
	assert.Contains(t, consoleOutput, "INF This is an info message")
	assert.Contains(t, consoleOutput, "WRN This is a warn message")
	assert.Contains(t, consoleOutput, "ERR This is an error message")
}

func TestNewHcLogAdapter_GetLevel(t *testing.T) {
	logger := NewLogger(
		context.Background(),
		LoggerConfig{
			Output:     []config.LogOutput{config.Console},
			Level:      zerolog.TraceLevel,
			TimeFormat: zerolog.TimeFormatUnix,
			NoColor:    true,
		},
	)

	hcLogAdapter := NewHcLogAdapter(&logger, "test")
	hcLogAdapter.SetLevel(hclog.Trace)
	assert.Equal(t, hclog.Trace, hcLogAdapter.GetLevel())

	hcLogAdapter.SetLevel(hclog.Debug)
	assert.Equal(t, hclog.Debug, hcLogAdapter.GetLevel())
	assert.NotEqual(t, zerolog.DebugLevel, logger.GetLevel(),
		"The logger should not be affected by the hclog adapter's level")
}
