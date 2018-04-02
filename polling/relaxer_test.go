package polling

import (
	"reflect"
	"testing"
	"time"

	"github.com/voicera/tester/assert"
)

func TestRelax_shouldNotRelax(t *testing.T) {
	relaxer := cyclicExponentialBackoffRelaxer{
		relaxationCondition: &passthrough{},
		sleep:               func(time.Duration) { t.Error("there's no rest for the wicked!") },
	}
	relaxer.relax(false)
}

func TestRelax_shouldRelax(t *testing.T) {
	shifter := uint(0)
	actualCallCount := 0
	expectedCallCount := 10
	relaxer := cyclicExponentialBackoffRelaxer{
		relaxationCondition:    &passthrough{},
		initialBackoffDuration: 1,
		currentBackoffDuration: 1,
		exponentialBackoffCap:  1000,
		sleep: func(d time.Duration) {
			assert.For(t).ThatActual(d).Equals(time.Duration(1 << shifter))
			shifter = (shifter + 1) % 10
			actualCallCount++
		},
	}

	for i := 0; i < expectedCallCount; i++ {
		relaxer.relax(true)
	}

	assert.For(t).ThatActual(actualCallCount).Equals(expectedCallCount)
}

func TestHooksAreHidden(t *testing.T) {
	assert.For(t).ThatType(reflect.TypeOf(cyclicExponentialBackoffRelaxer{})).HidesTestHooks()
}
