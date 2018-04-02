package validate

// NewValidationError creates a new validation error that can be handled
// differently in recovery. The parameter argumentName is optional.
func NewValidationError(err error, argumentName string) *ValidationError {
	return &ValidationError{error: err, ArgumentName: argumentName}
}

// ValidationError represents a validation error that can be handled differently
// in recovery.
type ValidationError struct {
	error
	ArgumentName string
}

func (err *ValidationError) Error() (message string) {
	message = err.error.Error()
	if len(err.ArgumentName) > 0 {
		message += "\nArgument: " + err.ArgumentName
	}
	return
}
