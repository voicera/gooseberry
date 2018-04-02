package zap

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Level is log level.
	Level = zap.DebugLevel

	// DefaultLogger is the default zap logger instance.
	DefaultLogger = newLoggerAdapter(false)

	inDevelopmentMode  = strings.EqualFold(os.Getenv("GOOSEBERRY_DEV_LOG"), "true")
	outputTestHook     = &bytes.Buffer{}
	encodeTimeFunction = encodeTimeUTC
)

// Logger implements log.LeveledLogger.
type Logger struct {
	*zap.Logger
}

func newLoggerAdapter(inUnitTestModeForThisPackage bool) *Logger {
	encoder := createEncoder()
	if inUnitTestModeForThisPackage {
		return &Logger{zap.New(zapcore.NewCore(encoder, zapcore.AddSync(outputTestHook), Level))}
	}
	core := zapcore.NewTee(
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stderr),
			zap.LevelEnablerFunc(func(x zapcore.Level) bool { return x >= zapcore.ErrorLevel })),
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(func(x zapcore.Level) bool { return x >= Level && x < zapcore.ErrorLevel })),
	)
	return &Logger{zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.ErrorLevel))}
}

func createEncoder() zapcore.Encoder {
	if flag.Lookup("test.v") != nil { // for any unit test, not just this package's, set the timestamp to 0
		encodeTimeFunction = func(_ time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendInt(0) }
	}

	if inDevelopmentMode {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = encodeTimeFunction
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = encodeTimeFunction
	return zapcore.NewJSONEncoder(encoderConfig)
}

func encodeTimeUTC(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(time.RFC3339))
}

// IsDebugEnabled checks whether or not debug logging is enabled.
// Useful to check before logging a message that requires preprocessing.
func IsDebugEnabled() bool {
	return DefaultLogger.IsDebugEnabled()
}

// Debug logs a debug message.
// It accepts varargs of alternating key and value parameters.
func Debug(message string, args ...interface{}) {
	DefaultLogger.Debug(message, args...)
}

// Info logs an informational message.
// It accepts varargs of alternating key and value parameters.
func Info(message string, args ...interface{}) {
	DefaultLogger.Info(message, args...)
}

// Warn logs a warning message.
// It accepts varargs of alternating key and value parameters.
func Warn(message string, args ...interface{}) {
	DefaultLogger.Warn(message, args...)
}

// Error logs an error message.
// It accepts varargs of alternating key and value parameters.
func Error(message string, args ...interface{}) {
	DefaultLogger.Error(message, args...)
}

// Sync flushes the Logger.
func Sync() error {
	return DefaultLogger.Sync()
}

// IsDebugEnabled checks whether or not debug logging is enabled.
// Useful to check before logging a message that requires preprocessing.
func (adapter *Logger) IsDebugEnabled() bool {
	return Level == zap.DebugLevel
}

// Debug logs a debug message.
// It accepts varargs of alternating key and value parameters.
func (adapter *Logger) Debug(message string, args ...interface{}) {
	adapter.Logger.Debug(message, convertLogParametersToZapFields(args)...)
}

// Info logs an informational message.
// It accepts varargs of alternating key and value parameters.
func (adapter *Logger) Info(message string, args ...interface{}) {
	adapter.Logger.Info(message, convertLogParametersToZapFields(args)...)
}

// Warn logs a warning message.
// It accepts varargs of alternating key and value parameters.
func (adapter *Logger) Warn(message string, args ...interface{}) {
	adapter.Logger.Warn(message, convertLogParametersToZapFields(args)...)
}

// Error logs an error message.
// It accepts varargs of alternating key and value parameters.
func (adapter *Logger) Error(message string, args ...interface{}) {
	adapter.Logger.Error(message, convertLogParametersToZapFields(args)...)
}

// convertLogParameters converts generic key value pairs into zap fields.
// TODO (Geish): move to SugaredLogger instead? P.S. this package predates SugaredLogger!
func convertLogParametersToZapFields(args []interface{}) []zapcore.Field {
	fields := []zapcore.Field{}
	for i := 0; i < len(args); i += 2 {
		key := args[i].(string)
		value := args[i+1]
		reflectedValue := reflect.ValueOf(value)
		if reflectedValue.Kind() == reflect.Ptr && reflectedValue.IsNil() {
			value = nil
		}
		fields = append(fields, zap.Any(key, value))
	}
	return fields
}
