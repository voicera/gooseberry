package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	jsonContentType = "application/json"
)

type jsonRequestCreator struct{}

func (jsonRequestCreator) CreateRequest(method string, url string, body interface{}) (*http.Request, error) {
	jsonBody, err := newJSONRequestBodyReader(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, url, jsonBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set(contentTypeHeaderKey, jsonContentType)
	return request, nil
}

func newJSONRequestBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	marshalled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(marshalled), nil
}

type jsonResponseDecoder struct{}

func (jsonResponseDecoder) DecodeResponse(body io.ReadCloser, result interface{}) error {
	return json.NewDecoder(body).Decode(result)
}
