package testutil

import (
	"fmt"
	"os"
)

// Logger is a logger for test purposes only.
type Logger struct {
	InDebugMode bool
}

// IsDebugEnabled checks whether or not debug logging is enabled.
func (logger *Logger) IsDebugEnabled() bool {
	return logger.InDebugMode
}

// Debug logs a debug message.
func (logger *Logger) Debug(message string, args ...interface{}) {
	if logger.IsDebugEnabled() {
		print("Debug", message, args)
	}
}

// Info logs a informational message.
func (*Logger) Info(message string, args ...interface{}) {
	print("Info", message, args)
}

// Warn logs a warning message.
func (*Logger) Warn(message string, args ...interface{}) {
	print("Warn", message, args)
}

// Error logs an error message.
func (*Logger) Error(message string, args ...interface{}) {
	print("Error", message, args)
}

// Sync flushes the logger.
func (*Logger) Sync() error {
	return os.Stdout.Sync()
}

func print(prefix string, message string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf("[%s] %s: ", prefix, message)
		fmt.Println(args...)
	} else {
		fmt.Printf("[%s] %s\n", prefix, message)
	}
}

// NewLogCapturer creates a logger that silently captures logs
// (for test purposes only).
func NewLogCapturer(inDebugMode bool) *LogCapturer {
	return &LogCapturer{
		Logger:        &Logger{InDebugMode: inDebugMode},
		DebugCaptures: []*CapturedLogEntry{},
		InfoCaptures:  []*CapturedLogEntry{},
		WarnCaptures:  []*CapturedLogEntry{},
		ErrorCaptures: []*CapturedLogEntry{},
	}
}

// LogCapturer is a logger that silently captures logs (for test purposes only).
type LogCapturer struct {
	*Logger
	DebugCaptures []*CapturedLogEntry
	InfoCaptures  []*CapturedLogEntry
	WarnCaptures  []*CapturedLogEntry
	ErrorCaptures []*CapturedLogEntry
}

// CapturedLogEntry represents a captured log entry (for test purposes only).
type CapturedLogEntry struct {
	Message   string
	Arguments []interface{}
}

// Debug captures a debug message.
func (logger *LogCapturer) Debug(message string, args ...interface{}) {
	if logger.IsDebugEnabled() {
		logger.DebugCaptures = append(logger.DebugCaptures, &CapturedLogEntry{message, args})
	}
}

// Info captures a informational message.
func (logger *LogCapturer) Info(message string, args ...interface{}) {
	logger.InfoCaptures = append(logger.InfoCaptures, &CapturedLogEntry{message, args})
}

// Warn captures a warning message.
func (logger *LogCapturer) Warn(message string, args ...interface{}) {
	logger.WarnCaptures = append(logger.WarnCaptures, &CapturedLogEntry{message, args})
}

// Error captures an error message.
func (logger *LogCapturer) Error(message string, args ...interface{}) {
	logger.ErrorCaptures = append(logger.ErrorCaptures, &CapturedLogEntry{message, args})
}
