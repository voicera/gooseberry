package validate

import (
	"errors"
	"testing"

	"github.com/voicera/tester/assert"
)

func TestValidationError_Error(t *testing.T) {
	cases := []struct {
		id                   string
		err                  *ValidationError
		expectedErrorMessage string
	}{
		{"empty argument name", NewValidationError(errors.New("42"), ""), "42"},
		{"unspecified argument name", &ValidationError{error: errors.New("42")}, "42"},
		{"non-empty argument name", NewValidationError(errors.New("42"), "answer"), "42\nArgument: answer"},
	}

	for _, c := range cases {
		assert.For(t, c.id).ThatActual(c.err.Error()).Equals(c.expectedErrorMessage)
	}
}
