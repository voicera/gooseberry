package log

import (
	"fmt"

	"github.com/voicera/gooseberry/testutil"
)

func ExampleNoOpLogger() {
	logger := &NoOpLogger{}
	fmt.Println("IsDebugEnabled:", logger.IsDebugEnabled())
	logger.Debug("Debug", "answer", 42)
	logger.Info("Info", "answer", 42)
	logger.Warn("Warn", "answer", 42)
	logger.Error("Error", "answer", 42)
	// Output:
	// IsDebugEnabled: false
}

func ExampleLeveledLogger_withPrefix() {
	logger := NewPrefixedLeveledLogger(&testutil.Logger{InDebugMode: true}, "prefix:")
	fmt.Println("IsDebugEnabled:", logger.IsDebugEnabled())
	logger.Debug("message", "answer", 42)
	logger.Info("message", "answer", 42)
	logger.Warn("message", "answer", 42)
	logger.Error("message", "answer", 42)
	// Output:
	// IsDebugEnabled: true
	// [Debug] prefix:message: [answer 42]
	// [Info] prefix:message: [answer 42]
	// [Warn] prefix:message: [answer 42]
	// [Error] prefix:message: [answer 42]
}
