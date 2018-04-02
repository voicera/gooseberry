package zap

import (
	"testing"
	"time"

	"github.com/voicera/tester/assert"
)

func TestProductionLogger(t *testing.T) {
	cases := []struct {
		logger        func(string, ...interface{})
		expectedLevel string
	}{
		{Debug, "debug"},
		{Info, "info"},
		{Warn, "warn"},
		{Error, "error"},
	}

	for _, c := range cases {
		setup(false)
		expected := `{"level":"` + c.expectedLevel + `","ts":0,"msg":"message","answer":42}` + "\n"
		c.logger("message", "answer", 42)
		assert.For(t, c.expectedLevel).ThatActualString(outputTestHook.String()).Equals(expected)
	}
}

func TestHumanReadableLogger(t *testing.T) {
	setup(true)
	Debug("message", "answer", 42)
	assert.For(t).ThatActualString(outputTestHook.String()).Equals("0\t\x1b[35mdebug\x1b[0m\tmessage\t{\"answer\": 42}\n")
}

func TestIsDebugEnabled_whenEnabled(t *testing.T) {
	assert.For(t).ThatActual(IsDebugEnabled()).IsTrue()
}

func TestNilValue(t *testing.T) {
	setup(false)
	var timestamp *time.Time
	Debug("message", "key", timestamp)
	expected := `{"level":"debug","ts":0,"msg":"message","key":null}` + "\n"
	assert.For(t).ThatActualString(outputTestHook.String()).Equals(expected)
}

func setup(inDevMode bool) {
	outputTestHook.Reset()
	inDevelopmentMode = inDevMode
	DefaultLogger = newLoggerAdapter(true)
}
