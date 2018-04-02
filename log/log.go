package log

// LeveledLogger provides a logging interface that supports structured logging.
// For example:
//     logger.Info("Found an answer.", "answer", 42)
//     logger.Debug("message", "key", value, "another key", anotherValue)
type LeveledLogger interface {
	// IsDebugEnabled checks whether or not debug logging is enabled.
	// Useful to check before logging a message that requires preprocessing.
	IsDebugEnabled() bool

	// Debug logs a debug message.
	// It accepts varargs of alternating key and value parameters.
	Debug(message string, args ...interface{})

	// Info logs a informational message.
	// It accepts varargs of alternating key and value parameters.
	Info(message string, args ...interface{})

	// Warn logs a warning message.
	// It accepts varargs of alternating key and value parameters.
	Warn(message string, args ...interface{})

	// Error logs an error message.
	// It accepts varargs of alternating key and value parameters.
	Error(message string, args ...interface{})

	// Sync flushes the logger.
	Sync() error
}

// NoOpLogger provides a NOOP logger (the default logging behavior).
type NoOpLogger struct{}

// IsDebugEnabled returns false.
func (*NoOpLogger) IsDebugEnabled() bool {
	return false
}

// Debug is a NOOP.
func (*NoOpLogger) Debug(string, ...interface{}) {}

// Info is a NOOP.
func (*NoOpLogger) Info(string, ...interface{}) {}

// Warn is a NOOP.
func (*NoOpLogger) Warn(string, ...interface{}) {}

// Error is a NOOP.
func (*NoOpLogger) Error(string, ...interface{}) {}

// Sync is a NOOP.
func (*NoOpLogger) Sync() error { return nil }

// NoOpLogger provides a NOOP logger (the default logging behavior).
type prefixedLogger struct {
	LeveledLogger
	prefix string
}

// NewPrefixedLeveledLogger prepends a prefix to log messages logged by
// the specified logger.
func NewPrefixedLeveledLogger(logger LeveledLogger, prefix string) LeveledLogger {
	return &prefixedLogger{LeveledLogger: logger, prefix: prefix}
}

func (logger *prefixedLogger) IsDebugEnabled() bool {
	return logger.LeveledLogger.IsDebugEnabled()
}

func (logger *prefixedLogger) Debug(message string, args ...interface{}) {
	logger.LeveledLogger.Debug(logger.prefix+message, args...)
}

func (logger *prefixedLogger) Info(message string, args ...interface{}) {
	logger.LeveledLogger.Info(logger.prefix+message, args...)
}

func (logger *prefixedLogger) Warn(message string, args ...interface{}) {
	logger.LeveledLogger.Warn(logger.prefix+message, args...)
}

func (logger *prefixedLogger) Error(message string, args ...interface{}) {
	logger.LeveledLogger.Error(logger.prefix+message, args...)
}

func (logger *prefixedLogger) Sync() error {
	return logger.LeveledLogger.Sync()
}
