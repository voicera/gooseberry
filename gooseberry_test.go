package gooseberry

import (
	"fmt"

	"github.com/voicera/gooseberry/testutil"
)

func Example_pluggableLogging() {
	originalLogger := Logger
	defer func() { Logger = originalLogger }()

	Logger = &testutil.Logger{InDebugMode: true}
	fmt.Printf("Logger: %T, IsDebugEnabled: %t\n", Logger, Logger.IsDebugEnabled())
	for _, log := range []func(string, ...interface{}){Logger.Debug, Logger.Info, Logger.Warn, Logger.Error} {
		log("message", "answer", 42)
	}
	// Output:
	// Logger: *testutil.Logger, IsDebugEnabled: true
	// [Debug] message: [answer 42]
	// [Info] message: [answer 42]
	// [Warn] message: [answer 42]
	// [Error] message: [answer 42]
}
