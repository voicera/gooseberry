package validate

import (
	"fmt"
)

const (
	inRangeErrorFormat = "value %v is out of range [%v, %v]."
)

// InRange validates that the specified value is in the specified range;
// argumentName is optional.
func InRange(value, inclusiveStart, inclusiveEnd float64, argumentName string) error {
	if value < inclusiveStart || value > inclusiveEnd {
		return &ValidationError{
			error:        fmt.Errorf(inRangeErrorFormat, value, inclusiveStart, inclusiveEnd),
			ArgumentName: argumentName,
		}
	}
	return nil
}
