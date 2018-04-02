package validate

import (
	"fmt"
	"testing"

	"github.com/voicera/tester/assert"
)

type rangeCheck struct {
	value      float64
	rangeStart float64
	rangeEnd   float64
}

func TestInRange_panic(t *testing.T) {
	cases := []rangeCheck{
		{4.2, 5, 42},
		{4.2, 0, 3.14},
	}

	for _, c := range cases {
		err := InRange(c.value, c.rangeStart, c.rangeEnd, "")
		expected := fmt.Sprintf("value %v is out of range [%v, %v].", c.value, c.rangeStart, c.rangeEnd)
		assert.For(t).ThatActualString(err.Error()).Equals(expected)
	}
}

func TestInRange_positive(t *testing.T) {
	cases := []rangeCheck{
		{1, 0, 2},
		{13, 13, 13},
		{3, -3, 13},
		{-3, -13, -3},
	}

	for _, c := range cases {
		err := InRange(c.value, c.rangeStart, c.rangeEnd, "argument")
		assert.For(t).ThatActual(err).IsNil()
	}
}
