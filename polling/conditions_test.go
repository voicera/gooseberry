package polling

import (
	"math/rand"
	"testing"

	"github.com/voicera/tester/assert"
)

type passthrough struct{}

func (*passthrough) shouldRelax(lastReceiptSucceeded bool) bool {
	return lastReceiptSucceeded
}

type not struct{}

func (*not) shouldRelax(lastReceiptSucceeded bool) bool {
	return !lastReceiptSucceeded
}

func TestAnd(t *testing.T) {
	cases := []struct {
		*and
		lastReceiptSucceeded bool
		expected             bool
	}{
		{&and{}, true, true},
		{&and{}, false, true},
		{&and{[]relaxationCondition{}}, true, true},
		{&and{[]relaxationCondition{}}, false, true},
		{&and{[]relaxationCondition{&and{}}}, true, true},
		{&and{[]relaxationCondition{&and{}}}, false, true},
		{&and{[]relaxationCondition{&passthrough{}}}, true, true},
		{&and{[]relaxationCondition{&passthrough{}}}, false, false},
		{&and{[]relaxationCondition{&passthrough{}, &not{}}}, true, false},
		{&and{[]relaxationCondition{&passthrough{}, &not{}}}, false, false},
	}

	for _, c := range cases {
		assert.For(t).ThatActual(c.shouldRelax(c.lastReceiptSucceeded)).Equals(c.expected)
	}
}

func TestEmptyHanded(t *testing.T) {
	condition := &emptyHanded{}
	assert.For(t).ThatActual(condition.shouldRelax(true)).IsFalse()
	assert.For(t).ThatActual(condition.shouldRelax(false)).IsTrue()
}

func TestBernoulliSampler(t *testing.T) {
	doNotCare := false
	condition, err := newBernoulliSampler(0.5)
	assert.For(t).ThatActual(err).IsNil()
	rand.Seed(0) // to produce the same sequence of pseudo-random numbers every time
	sequence := []bool{false, true, false, true, true, true, true, false, false, true}
	for _, expected := range sequence {
		assert.For(t).ThatActual(condition.shouldRelax(doNotCare)).Equals(expected)
	}
}
