package polling

import "time"

// Relaxer optionally causes pollers to relax (instead of busy-waiting)
// between polls.
type relaxer interface {
	// relax optionally causes pollers to relax.
	relax(lastReceiptSucceeded bool)
}

type cyclicExponentialBackoffRelaxer struct {
	relaxationCondition
	initialBackoffDuration time.Duration
	currentBackoffDuration time.Duration
	exponentialBackoffCap  time.Duration
	sleep                  func(time.Duration) `test-hook:"verify-unexported"`
}

func (relaxer *cyclicExponentialBackoffRelaxer) relax(lastReceiptSucceeded bool) {
	if !relaxer.shouldRelax(lastReceiptSucceeded) {
		return
	}

	relaxer.sleep(relaxer.currentBackoffDuration)
	relaxer.currentBackoffDuration += relaxer.currentBackoffDuration
	if relaxer.currentBackoffDuration > relaxer.exponentialBackoffCap {
		relaxer.currentBackoffDuration = relaxer.initialBackoffDuration
	}
}
