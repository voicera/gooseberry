// Package must provides functions used in package and variables initialization.
// If an error occurs, those functions will panic. Use judiciously!
package must

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/voicera/gooseberry/errors"
)

// Close wraps io.Closer.Close().
func Close(c io.Closer) {
	err := c.Close()
	errors.PanicIfNotNil(err)
}

// Getenv wraps os.Getenv(key).
func Getenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " is not set.")
	}

	return value
}

// LoadLocation wraps time.LoadLocation(name).
func LoadLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	errors.PanicIfNotNil(err)
	return location
}

// ParseBool wraps strconv.ParseBool(s).
func ParseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	errors.PanicIfNotNil(err)
	return b
}

// ParseDuration wraps time.ParseDuration(s).
func ParseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	errors.PanicIfNotNil(err)
	return duration
}

// ParseFloat wraps strconv.ParseFloat(s, 64).
func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	errors.PanicIfNotNil(err)
	return f
}

// ParseInt wraps strconv.ParseInt(s, 0, 64).
func ParseInt(s string) int64 {
	i, err := strconv.ParseInt(s, 0, 64)
	errors.PanicIfNotNil(err)
	return i
}

// ConvertToInt wraps strconv.Atoi for cases that ParseInt cannot handle
func ConvertToInt(s string) int {
	i, err := strconv.Atoi(s)
	errors.PanicIfNotNil(err)
	return i
}

// RemoveFileSystemNode wraps os.Remove(name).
func RemoveFileSystemNode(name string) {
	err := os.Remove(name)
	errors.PanicIfNotNil(err)
}
