package rest

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	urlEncodedContentType = "application/x-www-form-urlencoded"
)

type urlEncodedRequestCreator struct{}

func (urlEncodedRequestCreator) CreateRequest(method string, url string, body interface{}) (*http.Request, error) {
	return createRequest(method, url, body, newURLEncodedRequestBodyReader, urlEncodedContentType)
}

func newURLEncodedRequestBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return strings.NewReader(""), nil
	}
	parameters, ok := body.(map[string]string)
	if !ok {
		return nil, errors.New("body has to be map[string]string")
	}
	return newURLEncodedRequestReader(parameters), nil
}

func newURLEncodedRequestReader(parameters map[string]string) io.Reader {
	requestBody := url.Values{}
	for key, value := range parameters {
		requestBody.Add(key, value)
	}
	return strings.NewReader(requestBody.Encode())
}
