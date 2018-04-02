package web

import (
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/voicera/gooseberry/log"
)

const (
	strippedOutHeaderValue = "*******STRIPPED OUT*******"
	tokenReplacement       = `token":"` + strippedOutHeaderValue + `"`
)

var (
	tokenReplacer           = regexp.MustCompile(`token":".*?"`)
	sensitiveDataHeaderKeys = []string{"Authorization", "Cookie"}
)

// NewBasicAuthRoundTripper creates a RoundTripper that decorates another
// round tripper by adding basic auth using the specified username and password.
func NewBasicAuthRoundTripper(roundTripper http.RoundTripper, username, password string) http.RoundTripper {
	return &basicAuthRoundTripper{innerRoundTripper: roundTripper, username: username, password: password}
}

type basicAuthRoundTripper struct {
	innerRoundTripper http.RoundTripper
	username          string
	password          string
}

func (roundTripper *basicAuthRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth(roundTripper.username, roundTripper.password)
	return roundTripper.innerRoundTripper.RoundTrip(request)
}

// NewCustomHeadersRoundTripper creates a RoundTripper that decorates another
// round tripper by adding custom headers.
func NewCustomHeadersRoundTripper(roundTripper http.RoundTripper, headers map[string]string) http.RoundTripper {
	return &customHeadersRoundTripper{innerRoundTripper: roundTripper, headers: headers}
}

type customHeadersRoundTripper struct {
	innerRoundTripper http.RoundTripper
	headers           map[string]string
}

func (roundTripper *customHeadersRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	for key, value := range roundTripper.headers {
		request.Header.Add(key, value)
	}
	return roundTripper.innerRoundTripper.RoundTrip(request)
}

// NewLeveledLoggerRoundTripper creates a RoundTripper that decorates another
// round tripper by using the specified leveled logger to log requests
// and responses only in the debug log level; and errors in all logging levels.
func NewLeveledLoggerRoundTripper(roundTripper http.RoundTripper, logger log.LeveledLogger) http.RoundTripper {
	return &loggingRoundTripper{innerRoundTripper: roundTripper, logger: logger}
}

type loggingRoundTripper struct {
	innerRoundTripper http.RoundTripper
	logger            log.LeveledLogger
}

func (roundTripper *loggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	roundTripper.debugLogRequest(request)
	response, err := roundTripper.innerRoundTripper.RoundTrip(request)
	roundTripper.debugLogResponse(response, err)
	return response, err
}

func (roundTripper *loggingRoundTripper) debugLogRequest(request *http.Request) {
	if roundTripper.logger.IsDebugEnabled() {
		censoredHeaders := map[string]string{}
		for _, headerKey := range sensitiveDataHeaderKeys {
			headerValue := request.Header.Get(headerKey)
			if headerValue != "" {
				request.Header.Set(headerKey, strippedOutHeaderValue)
				censoredHeaders[headerKey] = headerValue
			}
		}
		dump, err := httputil.DumpRequestOut(request, true)
		for key, value := range censoredHeaders { // restore censored headers
			request.Header.Set(key, value)
		}
		if err == nil {
			roundTripper.logger.Debug("Request", "request", string(dump))
		}
	}
}

func (roundTripper *loggingRoundTripper) debugLogResponse(response *http.Response, responseError error) {
	if responseError != nil {
		if response != nil {
			dump, err := httputil.DumpResponse(response, true)
			if err == nil {
				roundTripper.logger.Error("Response error", "responseError", responseError, "response", string(dump))
			}
		} else {
			roundTripper.logger.Error("Response error", "responseError", responseError)
		}
	} else if roundTripper.logger.IsDebugEnabled() {
		dump, err := httputil.DumpResponse(response, true)
		if err == nil {
			roundTripper.logger.Debug("Response", "response", stripOutSensitiveData(string(dump)))
		}
	}
}

func stripOutSensitiveData(s string) string {
	return tokenReplacer.ReplaceAllString(s, tokenReplacement)
}
