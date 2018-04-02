package web

import (
	"fmt"
)

// HTTPError represents an HTTP error that can be handled differently
// in recovery.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("HTTP Status Code %d: %s", err.StatusCode, err.Body)
}
