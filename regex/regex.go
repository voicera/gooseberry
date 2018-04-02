// Package regex provides regular expressions utilities.
package regex

import (
	"regexp"
)

// FindNamedStringSubmatch returns a map whose elements are the names of
// explicitly named capturing groups in the specified pattern (denoted by
// ?P<name>) and their respective submatches in the input string, as
// defined by the 'Submatch' description in the regexp package comment.
// A return value of nil indicates no match (to be consistent with regexp).
// A return value of an empty map indicates lack of any named capturing group.
func FindNamedStringSubmatch(pattern *regexp.Regexp, input string) map[string]string {
	matches := pattern.FindStringSubmatch(input)
	if matches == nil {
		return nil
	}

	names := pattern.SubexpNames()
	mapping := make(map[string]string, len(names))
	for i, key := range names {
		if key != "" {
			mapping[key] = matches[i]
		}
	}

	return mapping
}
