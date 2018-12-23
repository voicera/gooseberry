// Package errors provides utilities for error handling.
package errors

import (
	"bytes"
	"fmt"
	"sync"
)

var (
	bufferPool = sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}
)

// ErrorString is a trivial implementation of error.
type ErrorString string

// AggregateError represents an aggregate error.
type AggregateError struct {
	Header string
	Errors []error
}

func (err ErrorString) Error() string {
	return string(err)
}

// NewAggregateError creates a new aggregate error with the specified header.
func NewAggregateError(header string, errors ...error) *AggregateError {
	return &AggregateError{Header: header, Errors: errors}
}

func (aggregateError *AggregateError) Error() string {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	if aggregateError.Header != "" {
		if _, err := buffer.WriteString(aggregateError.Header); err != nil {
			return fmt.Sprint(aggregateError)
		}
		if _, err := buffer.WriteRune('\n'); err != nil {
			return fmt.Sprint(aggregateError)
		}
	}
	for _, e := range aggregateError.Errors {
		if _, err := buffer.WriteString(e.Error()); err != nil {
			return fmt.Sprint(aggregateError)
		}
		if _, err := buffer.WriteRune('\n'); err != nil {
			return fmt.Sprint(aggregateError)
		}
	}
	message := buffer.String()
	bufferPool.Put(buffer)
	return message
}

// PanicIfNotNil panics if the specified error is not nil.
func PanicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
