package web

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/voicera/gooseberry/testutil"
	"github.com/voicera/tester/assert"
)

const pingRequestDump = "GET / HTTP/1.1\r\n" +
	"Host: host\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 5\r\nAccept-Encoding: gzip\r\n\r\nping!"

func TestBasicAuthRoundTripper(t *testing.T) {
	request, _ := http.NewRequest("GET", "//host", strings.NewReader(""))
	innerRoundTripper := &mockRoundTripper{}
	roundTripper := NewBasicAuthRoundTripper(innerRoundTripper, "username", "password")
	_, err := roundTripper.RoundTrip(request)
	if assert.For(t).ThatActual(err).IsNil().Passed() {
		username, password, _ := innerRoundTripper.request.BasicAuth()
		assert.For(t).ThatActualString(username).Equals("username")
		assert.For(t).ThatActualString(password).Equals("password")
	}
}

func TestCustomHeadersRoundTripper(t *testing.T) {
	request, _ := http.NewRequest("GET", "//host", strings.NewReader(""))
	innerRoundTripper := &mockRoundTripper{}

	cases := []struct {
		id       string
		headers  map[string]string
		expected map[string]string
	}{
		{"nil map", nil, map[string]string{"foo": "", "bar": ""}},
		{"empty map", map[string]string{}, map[string]string{"foo": "", "bar": ""}},
		{"vanilla", map[string]string{"foo": "oof", "bar": "rab"}, map[string]string{"foo": "oof", "bar": "rab"}},
	}

	for _, c := range cases {
		roundTripper := NewCustomHeadersRoundTripper(innerRoundTripper, c.headers)
		_, err := roundTripper.RoundTrip(request)
		if assert.For(t).ThatActual(err).IsNil().Passed() {
			for key, expected := range c.expected {
				actual := innerRoundTripper.request.Header.Get(key)
				assert.For(t, c.id, key).ThatActualString(actual).Equals(expected)
			}
		}
	}
}

func TestLeveledLoggerRoundTripper_vanilla(t *testing.T) {
	response := &http.Response{Body: ioutil.NopCloser(strings.NewReader("ack!"))}
	cases := []struct {
		id          string
		inDebugMode bool
		expected    []*testutil.CapturedLogEntry
	}{
		{"debug logging disabled", false, []*testutil.CapturedLogEntry{}},
		{"debug logging enabled", true, []*testutil.CapturedLogEntry{
			{"Request", []interface{}{"request", pingRequestDump}},
			{"Response", []interface{}{"response", "HTTP/0.0 000 status code 0\r\n\r\nack!"}},
		}},
	}

	for _, c := range cases {
		capturedLogs := testLoggingRoundTripper(t, c.inDebugMode, response, nil).DebugCaptures
		assert.For(t, c.id).ThatActual(capturedLogs).Equals(c.expected).ThenDiffOnFail()
	}
}

func TestLeveledLoggerRoundTripper_onError(t *testing.T) {
	expectedError := errors.New("to err is human")
	response := &http.Response{Body: ioutil.NopCloser(strings.NewReader("!"))}
	cases := []struct {
		id       string
		response *http.Response
		expected []*testutil.CapturedLogEntry
	}{
		{"no response", nil, []*testutil.CapturedLogEntry{
			{"Response error", []interface{}{"responseError", expectedError}},
		}},
		{"some response", response, []*testutil.CapturedLogEntry{
			{"Response error", []interface{}{"responseError", expectedError, "response", "HTTP/0.0 000 status code 0\r\n\r\n!"}},
		}},
	}

	for _, c := range cases {
		capturedLogs := testLoggingRoundTripper(t, true, c.response, expectedError).ErrorCaptures
		assert.For(t, c.id).ThatActual(capturedLogs).Equals(c.expected).ThenDiffOnFail()
	}
}

func testLoggingRoundTripper(
	t *testing.T, debug bool, mockResponse *http.Response, mockError error) *testutil.LogCapturer {
	logCapturer := testutil.NewLogCapturer(debug)
	request, _ := http.NewRequest("GET", "http://host", strings.NewReader("ping!"))
	roundTripper := NewLeveledLoggerRoundTripper(&mockRoundTripper{response: mockResponse, err: mockError}, logCapturer)
	_, err := roundTripper.RoundTrip(request)
	assert.For(t).ThatActual(err).Equals(mockError)
	return logCapturer
}

type mockRoundTripper struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (roundTripper *mockRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	roundTripper.request = request
	return roundTripper.response, roundTripper.err
}
