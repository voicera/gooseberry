// Package gooseberry provides common Go libraries to import in other Go projects.
package gooseberry

import (
	"github.com/voicera/gooseberry/log"
)

// Logger provides a hook for importing applications to wire up their own logger
// for gooseberry to use. By default, logging in gooseberry is a NOOP.
// Optionally, one can set this logger to zap.Logger from the log/zap package.
var Logger log.LeveledLogger = &log.NoOpLogger{}
