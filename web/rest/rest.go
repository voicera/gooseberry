package rest

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/voicera/gooseberry"
	"github.com/voicera/gooseberry/web"
)

const (
	contentTypeHeaderKey = "Content-Type"
	userAgentHeaderKey   = "User-Agent"
	userAgentHeaderValue = "gooseberry"
)

var (
	jsonRequestCreatorInstance       = jsonRequestCreator{}
	jsonResponseDecoderInstance      = jsonResponseDecoder{}
	noopResponseDecoderInstance      = noopResponseDecoder{}
	urlEncodedRequestCreatorInstance = urlEncodedRequestCreator{}
)

type bodyEncoder func(body interface{}) (io.Reader, error)

type requestCreator interface {
	CreateRequest(method string, url string, body interface{}) (*http.Request, error)
}

type responseDecoder interface {
	DecodeResponse(body io.ReadCloser, result interface{}) error
}

type client struct {
	requestCreator
	responseDecoder
	baseURL    string
	httpClient *http.Client
}

// NewJSONClient creates a new REST client that uses JSON to encode requests
// and decode responses.
func NewJSONClient(httpClient *http.Client) Client {
	return &client{
		httpClient:      httpClient,
		requestCreator:  jsonRequestCreatorInstance,
		responseDecoder: jsonResponseDecoderInstance,
	}
}

// NewURLEncodedRequestJSONResponseClient creates a new REST client that
// URL-encodes requests and uses JSON to decode responses. The body parameter
// for this client's methods must be a map[string]string instance.
func NewURLEncodedRequestJSONResponseClient(httpClient *http.Client) Client {
	return &client{
		httpClient:      httpClient,
		requestCreator:  urlEncodedRequestCreatorInstance,
		responseDecoder: jsonResponseDecoderInstance,
	}
}

// NewJSONRequestNoopResponseClient creates a new REST client that uses JSON
// to encode requests and does not decode the response (it leaves the decoding
// to the caller, which has access to the response object).
func NewJSONRequestNoopResponseClient(httpClient *http.Client) Client {
	return &client{
		httpClient:      httpClient,
		requestCreator:  jsonRequestCreatorInstance,
		responseDecoder: noopResponseDecoderInstance,
	}
}

// Client represts a web client to use with REST APIs.
type Client interface {
	// WithBaseURL configures the client with a base URL; returns modified self.
	WithBaseURL(baseURL string) Client

	// Get makes a GET request.
	Get(url string, body interface{}, result interface{}) (*http.Response, error)

	// Head makes a HEAD request.
	Head(url string, body interface{}, result interface{}) (*http.Response, error)

	// Post makes a POST request.
	Post(url string, body interface{}, result interface{}) (*http.Response, error)

	// Put makes a PUT request.
	Put(url string, body interface{}, result interface{}) (*http.Response, error)

	// Patch makes a PATCH request.
	Patch(url string, body interface{}, result interface{}) (*http.Response, error)

	// Delete makes a DELETE request.
	Delete(url string, body interface{}, result interface{}) (*http.Response, error)

	// Do makes a REST request using JSON for input and output.
	Do(method string, url string, body interface{}, result interface{}) (*http.Response, error)
}

func (c *client) WithBaseURL(baseURL string) Client {
	if strings.HasSuffix(baseURL, "/") {
		c.baseURL = baseURL
	} else {
		c.baseURL = baseURL + "/"
	}
	return c
}

func (c *client) Get(url string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Do(http.MethodGet, url, body, result)
}

func (c *client) Head(url string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Do(http.MethodHead, url, body, result)
}

func (c *client) Post(url string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Do(http.MethodPost, url, body, result)
}

func (c *client) Put(url string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Do(http.MethodPut, url, body, result)
}

func (c *client) Patch(url string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Do(http.MethodPatch, url, body, result)
}

func (c *client) Delete(url string, body interface{}, result interface{}) (*http.Response, error) {
	return c.Do(http.MethodDelete, url, body, result)
}

func (c *client) Do(
	method string, url string, body interface{}, result interface{}) (*http.Response, error) {
	request, err := c.CreateRequest(method, c.resolveURL(url), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set(userAgentHeaderKey, userAgentHeaderValue)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return response, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			gooseberry.Logger.Error("Error closing response", "status", response.Status, "request", request.URL)
		}
	}()

	if response.StatusCode/100 != 2 { // if not 2xx Success; must be handled here
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return response, err
		}
		// TODO (Geish): handle other status codes
		return nil, &web.HTTPError{StatusCode: response.StatusCode, Body: string(body)}
	}

	if result != nil {
		if err := c.DecodeResponse(response.Body, result); err != nil {
			return response, err
		}
	}
	return response, nil
}

func (c *client) resolveURL(path string) string {
	// TODO (Geish): if path is an absolute URL, take it and ignore the base one
	return c.baseURL + path
}

func createRequest(
	method string, url string, body interface{}, encode bodyEncoder, contentType string) (*http.Request, error) {
	bodyReader, err := encode(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set(contentTypeHeaderKey, contentType)
	return request, nil
}
