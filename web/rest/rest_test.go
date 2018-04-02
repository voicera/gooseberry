package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/voicera/gooseberry/web"
	"github.com/voicera/tester/assert"
)

const (
	restNounPrefix = "/rest"
)

var (
	expectedRequestBody = &map[string]string{"request": "body"}
	expectedResult      = &map[string]int{"answer": 42}
	currentPortNumber   = int32(1234)
)

func TestJSONClient_vanilla(t *testing.T) {
	cases := []*testCase{
		{"JSON_GetNilRequestBody", nil, http.MethodGet, expectedResult},
		{"JSON_GetNonNilRequestBody", expectedRequestBody, http.MethodGet, expectedResult},
		{"JSON_PostNilRequestBody", nil, http.MethodPost, expectedResult},
		{"JSON_PostNonNilRequestBody", expectedRequestBody, http.MethodPost, expectedResult},
		{"JSON_PutNilRequestBody", nil, http.MethodPut, expectedResult},
		{"JSON_PutNonNilRequestBody", expectedRequestBody, http.MethodPut, expectedResult},
		{"JSON_DeleteNilRequestBody", nil, http.MethodDelete, expectedResult},
		{"JSON_DeleteNonNilRequestBody", expectedRequestBody, http.MethodDelete, expectedResult},
		{"JSON_OptionsNilRequestBody", nil, http.MethodOptions, expectedResult},
		{"JSON_OptionsNonNilRequestBody", expectedRequestBody, http.MethodOptions, expectedResult},
	}

	for _, c := range cases {
		http.HandleFunc(restNounPrefix+c.id, func(writer http.ResponseWriter, request *http.Request) {
			validateJSONRequest(t, c, request)
			err := json.NewEncoder(writer).Encode(c.result)
			assert.For(t).ThatActual(err).IsNil()
		})

		listener, url := mustListen(t, c)
		defer func() {
			err := listener.Close()
			assert.For(t).ThatActual(err).IsNil()
		}()

		client := NewJSONClient(http.DefaultClient)
		result := &map[string]int{}
		var err error

		if c.requestMethod == http.MethodGet {
			_, err = client.Get(url, c.requestBody, result)
		} else if c.requestMethod == http.MethodPost {
			_, err = client.Post(url, c.requestBody, result)
		} else if c.requestMethod == http.MethodPut {
			_, err = client.Put(url, c.requestBody, result)
		} else if c.requestMethod == http.MethodDelete {
			_, err = client.Delete(url, c.requestBody, result)
		} else {
			_, err = client.Do(c.requestMethod, url, c.requestBody, result)
		}

		if assert.For(t, c.id).ThatActual(err).IsNil().Passed() {
			assert.For(t, c.id).ThatActual(result).Equals(c.result)
		}
	}
}

func TestJSONClient_onError(t *testing.T) {
	c := &testCase{"JSON_GetNotFound", nil, http.MethodGet, expectedResult}
	http.HandleFunc(restNounPrefix+c.id, func(writer http.ResponseWriter, request *http.Request) {
		validateJSONRequest(t, c, request)
		http.NotFound(writer, request)
	})

	listener, url := mustListen(t, c)
	defer func() {
		err := listener.Close()
		assert.For(t).ThatActual(err).IsNil()
	}()

	client := NewJSONClient(http.DefaultClient)
	result := &map[string]int{}
	_, err := client.Get(url, c.requestBody, result)
	if assert.For(t, c.id).ThatActual(err).IsNotNil().Passed() {
		httpError := err.(*web.HTTPError)
		assert.For(t, c.id).ThatActual(httpError.StatusCode).Equals(404)
		assert.For(t, c.id).ThatActualString(httpError.Body).Equals("404 page not found\n")
	}
}

func TestJSONResponseClient_vanilla(t *testing.T) {
	cases := []*testCase{
		{"JSONDecoder_Get", *expectedRequestBody, http.MethodGet, expectedResult},
		{"JSONDecoder_Post", *expectedRequestBody, http.MethodPost, expectedResult},
		{"JSONDecoder_Put", *expectedRequestBody, http.MethodPut, expectedResult},
		{"JSONDecoder_Delete", *expectedRequestBody, http.MethodDelete, expectedResult},
		{"JSONDecoder_Options", *expectedRequestBody, http.MethodOptions, expectedResult},
	}

	for _, c := range cases {
		http.HandleFunc(restNounPrefix+c.id, func(writer http.ResponseWriter, request *http.Request) {
			validateURLEncodedRequest(t, c, request)
			err := json.NewEncoder(writer).Encode(c.result)
			assert.For(t).ThatActual(err).IsNil()
		})

		listener, url := mustListen(t, c)
		defer func() {
			err := listener.Close()
			assert.For(t).ThatActual(err).IsNil()
		}()

		client := NewURLEncodedRequestJSONResponseClient(http.DefaultClient)
		result := &map[string]int{}
		requestBody := c.requestBody.(map[string]string)
		var err error

		if c.requestMethod == http.MethodGet {
			_, err = client.Get(url, requestBody, result)
		} else if c.requestMethod == http.MethodPost {
			_, err = client.Post(url, requestBody, result)
		} else if c.requestMethod == http.MethodPut {
			_, err = client.Put(url, requestBody, result)
		} else if c.requestMethod == http.MethodDelete {
			_, err = client.Delete(url, requestBody, result)
		} else {
			_, err = client.Do(c.requestMethod, url, requestBody, result)
		}

		if assert.For(t, c.id).ThatActual(err).IsNil().Passed() {
			assert.For(t, c.id).ThatActual(result).Equals(c.result)
		}
	}
}

func TestJSONResponseClient_onError(t *testing.T) {
	c := &testCase{"JSONDecoder_GetNotFound", expectedRequestBody, http.MethodGet, expectedResult}
	http.HandleFunc(restNounPrefix+c.id, func(writer http.ResponseWriter, request *http.Request) {
		validateURLEncodedRequest(t, c, request)
		http.NotFound(writer, request)
	})

	listener, url := mustListen(t, c)
	defer func() {
		err := listener.Close()
		assert.For(t).ThatActual(err).IsNil()
	}()

	client := NewURLEncodedRequestJSONResponseClient(http.DefaultClient)
	result := &map[string]int{}
	_, err := client.Get(url, *expectedRequestBody, result)
	if assert.For(t, c.id).ThatActual(err).IsNotNil().Passed() {
		httpError := err.(*web.HTTPError)
		assert.For(t, c.id).ThatActual(httpError.StatusCode).Equals(404)
		assert.For(t, c.id).ThatActualString(httpError.Body).Equals("404 page not found\n")
	}
}

func validateJSONRequest(t *testing.T, c *testCase, request *http.Request) {
	requestBody := &map[string]string{}
	defer func() {
		err := request.Body.Close()
		assert.For(t).ThatActual(err).IsNil()
	}()

	err := json.NewDecoder(request.Body).Decode(requestBody)
	if err != nil {
		assert.For(t, c.id).ThatActual(c.requestBody).IsNil()
		assert.For(t, c.id).ThatActualString(err.Error()).Equals("EOF")
	} else {
		assert.For(t, c.id).ThatActual(requestBody).Equals(c.requestBody)
	}

	validateRequest(t, c, request, "application/json")
}

func validateURLEncodedRequest(t *testing.T, c *testCase, request *http.Request) {
	defer func() {
		err := request.Body.Close()
		assert.For(t).ThatActual(err).IsNil()
	}()

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		assert.For(t, c.id).ThatActual(c.requestBody).IsNil()
		assert.For(t, c.id).ThatActualString(err.Error()).Equals("EOF")
	} else {
		urlValues, err := url.ParseQuery(string(requestBody))
		assert.For(t, c.id).ThatActual(err).IsNil()
		assert.For(t, c.id).ThatActualString(urlValues.Get("request")).Equals("body")
	}

	validateRequest(t, c, request, "application/x-www-form-urlencoded")
}

func validateRequest(t *testing.T, c *testCase, request *http.Request, expectedContentType string) {
	assert.For(t, c.id).ThatActualString(request.Header.Get("Content-Type")).Equals(expectedContentType)
	assert.For(t, c.id).ThatActualString(request.Header.Get("User-Agent")).Equals("gooseberry")
	assert.For(t, c.id).ThatActualString(request.Method).Equals(c.requestMethod)
}

func mustListen(t *testing.T, c *testCase) (net.Listener, string) {
	// Each test runs using a new port number in parallel
	serverAddress := fmt.Sprintf(":%d", atomic.AddInt32(&currentPortNumber, 1))
	url := "http://localhost" + serverAddress + restNounPrefix + c.id
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		t.Fatal(err)
	}

	go mustServe(listener)
	return listener, url
}

func mustServe(listener net.Listener) {
	if err := http.Serve(listener, nil); err != nil {
		log.Println(err)
	}
}

type testCase struct {
	id            string
	requestBody   interface{}
	requestMethod string
	result        interface{}
}
