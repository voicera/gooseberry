package regex

import (
	"regexp"
	"testing"

	"github.com/voicera/tester/assert"
)

func TestFindNamedStringSubmatch(t *testing.T) {
	cases := []struct {
		id             string
		pattern        *regexp.Regexp
		s              string
		expectedOutput map[string]string
	}{
		{"match - empty input", regexp.MustCompile(".*"), "", map[string]string{}},
		{"match - no subexpressions", regexp.MustCompile(".*"), "input", map[string]string{}},
		{"match - named subexpression", regexp.MustCompile(`1(?P<foo>\d)`), "13", map[string]string{"foo": "3"}},
		{"match - unnamed subexpressions", regexp.MustCompile(`(\d)(.)(x(\d))?`), "13", map[string]string{}},
		{
			"match - a mix of named and unnamed subexpressions",
			regexp.MustCompile(`(?P<x>.)(?P<y>.)(42)?`),
			"13",
			map[string]string{"x": "1", "y": "3"},
		},
		{"mismatch - empty input", regexp.MustCompile("x"), "", nil},
		{"mismatch - no subexpressions", regexp.MustCompile("x"), "input", nil},
		{"mismatch - named subexpression", regexp.MustCompile(`x(?P<foo>\d)`), "input", nil},
		{"mismatch - named subexpressions", regexp.MustCompile(`(?P<x>1)(?P<y>3)`), "input", nil},
		{"mismatch - unnamed subexpressions", regexp.MustCompile(`(\d)(1)(x(\d))?`), "input", nil},
		{
			"mismatch - a mix of named and unnamed subexpressions",
			regexp.MustCompile(`(?P<x>x)(?P<y>y)(42)?`),
			"13",
			nil,
		},
	}

	for _, c := range cases {
		assert.For(t, c.id).ThatActual(FindNamedStringSubmatch(c.pattern, c.s)).Equals(c.expectedOutput)
	}
}
