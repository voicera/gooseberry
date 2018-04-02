# Gooseberry: Common Packages for Go Microservices

[![Build Status](https://travis-ci.org/voicera/gooseberry.svg?branch=master)](https://travis-ci.org/voicera/gooseberry)
[![Go Report Card](https://goreportcard.com/badge/github.com/voicera/gooseberry)](https://goreportcard.com/report/github.com/voicera/gooseberry)
[![GoDoc](https://godoc.org/github.com/voicera/gooseberry?status.svg)](https://godoc.org/github.com/voicera/gooseberry)
[![Maintainability](https://api.codeclimate.com/v1/badges/98e8195b41246e1c573d/maintainability)](https://codeclimate.com/github/voicera/gooseberry/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/98e8195b41246e1c573d/test_coverage)](https://codeclimate.com/github/voicera/gooseberry/test_coverage)

Gooseberry is a collection of common Go packages that Voicera uses in microservices.
It's an incomplete library, named after a fruit that looks like an ungrown clementine.
We'd like to build gooseberry to be like Guava is for Java.

## Features
* REST clients, web client with logging, basic auth support, etc.
* Container structs like immutable maps, priority queues, sets, etc.
* Error aggregation (multiple errors into one with a header message)
* Leveled logger with a prefix and a wrapper for [zap](go.uber.org/zap)
* Polling with an exponential backoff and a Bernoulli trials for resetting
* Uniform Resource Name struct that implements [RFC8141](https://tools.ietf.org/html/rfc8141)

## Quick Start
To get the latest version: `go get -u github.com/voicera/gooseberry`

### Example
The example below creates a RESTful Twilio client to make a phone call and poll
for call history. The client uses a logger for requests and responses, keeps
polling for calls made using an exponential backoff poller.

```go
package main

import (
	"net/http"
	"time"

	"github.com/voicera/gooseberry"
	"github.com/voicera/gooseberry/log"
	"github.com/voicera/gooseberry/log/zap"
	"github.com/voicera/gooseberry/polling"
	"github.com/voicera/gooseberry/web"
	"github.com/voicera/gooseberry/web/rest"
)

const (
	baseURL    = "https://api.twilio.com/2010-04-01/Accounts/"
	accountSid = "AC072dcbab90350495b2c0fabf9a7817bb"
	authToken  = "883XXXXXXXXXXXXXXXXXXXXXXXXX1985"
)

type call struct {
	SID    string `json:"sid"`
	Status string `json:"status"`
}

type receiver struct {
	restClient rest.Client
}

func main() {
	gooseberry.Logger = zap.DefaultLogger
	gooseberry.Logger.Info("starting example")
	transport := web.NewBasicAuthRoundTripper(
		web.NewLeveledLoggerRoundTripper(
			http.DefaultTransport,
			log.NewPrefixedLeveledLogger(gooseberry.Logger, "TWL:")),
		accountSid, authToken)
	httpClient := &http.Client{Transport: transport}
	twilioClient := rest.NewURLEncodedRequestJSONResponseClient(httpClient).
		WithBaseURL(baseURL + accountSid)
	makeCall(twilioClient)
	poll(&receiver{twilioClient})
	time.Sleep(2 * time.Second)
	gooseberry.Logger.Sync()
}

func makeCall(twilioClient rest.Client) {
	parameters := map[string]string{
		"From": "+15005550006",
		"To":   "+14108675310",
		"Url":  "http://demo.twilio.com/docs/voice.xml",
	}
	call := &call{}
	if _, err := twilioClient.Post("Calls.json", parameters, &call); err != nil {
		gooseberry.Logger.Error("error making a call", "err", err)
	} else {
		gooseberry.Logger.Debug("made a call", "sid", call.SID)
	}
}

func poll(receiver *receiver) {
	poller, err := polling.NewBernoulliExponentialBackoffPoller(
		receiver, "twilio", 0.95, time.Second, time.Minute)
	if err != nil {
		gooseberry.Logger.Error("error creating a poller", "err", err)
	}
	go poller.Start()
}

func (r *receiver) Receive() (interface{}, bool, error) {
	calls := []*call{}
	_, err := r.restClient.Get("Calls", nil, &calls)
	return calls, len(calls) > 0, err
}
```

Running the above example produces an output that looks like the following (which was heavily edited for brevity):

```bash
{"level":"info","ts":"2018-04-02T00:58:05Z","caller":"runtime/proc.go:198","msg":"starting example"}
{"level":"debug","ts":"2018-04-02T00:58:05Z","caller":"web/web.go:90","msg":"TWL:Request","request":"POST /2010-04-01/Accounts/AC072dcbab90350495b2c0fabf9a7817bb/Calls.json HTTP/1.1\r\nHost: api.twilio.com\r\nUser-Agent: gooseberry\r\nContent-Length: 89\r\nAuthorization: *******STRIPPED OUT*******\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\nFrom=%2B15005550006&To=%2B14108675310&Url=http%3A%2F%2Fdemo.twilio.com%2Fdocs%2Fvoice.xml"}
{"level":"debug","ts":"2018-04-02T00:58:06Z","caller":"web/web.go:108","msg":"TWL:Response","response":"HTTP/1.1 401 UNAUTHORIZED
...
\"Your AccountSid or AuthToken was incorrect.\", \"message\": \"Authenticate\", \"more_info\": \"https://www.twilio.com/docs/errors/20003\", \"status\": 401}"}
{"level":"error","ts":"2018-04-02T00:58:06Z","caller":"example/main.go:40","msg":"error making a call","err":"HTTP Status Code 401
...
runtime.main\n\t/usr/local/go/src/runtime/proc.go:198"}
{"level":"debug","ts":"2018-04-02T00:58:06Z","caller":"runtime/asm_amd64.s:2361","msg":"Started","poller":"twilio"}
{"level":"debug","ts":"2018-04-02T00:58:06Z","caller":"web/web.go:90","msg":"TWL:Request","request":"GET /2010-04-01/Accounts/AC072dcbab90350495b2c0fabf9a7817bb/Calls HTTP/1.1\r\nHost: api.twilio.com\r\nUser-Agent: gooseberry\r\nAuthorization: *******STRIPPED OUT*******\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\n"}
{"level":"debug","ts":"2018-04-02T00:58:06Z","caller":"web/web.go:108","msg":"TWL:Response","response":"HTTP/1.1 401 UNAUTHORIZED\r\nContent-Length:
...
X-Shenanigans: none\r\n\r\n<?xml version='1.0' encoding='UTF-8'?>\n<TwilioResponse><RestException><Code>20003</Code><Detail>Your AccountSid or AuthToken was incorrect.</Detail><Message>Authenticate</Message><MoreInfo>https://www.twilio.com/docs/errors/20003</MoreInfo><Status>401</Status></RestException></TwilioResponse>"}
...
{"level":"debug","ts":"2018-04-02T00:58:06Z","caller":"polling/poller.go:98","msg":"Relaxing","poller":"twilio"}
{"level":"debug","ts":"2018-04-02T00:58:07Z","caller":"web/web.go:90","msg":"TWL:Request","request":"GET /2010-04-01/Accounts/AC072dcbab90350495b2c0fabf9a7817bb/Calls HTTP/1.1\r\nHost: api.twilio.com\r\nUser-Agent: gooseberry\r\nAuthorization: *******STRIPPED OUT*******\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\n"}
...
{"level":"debug","ts":"2018-04-02T00:58:07Z","caller":"polling/poller.go:98","msg":"Relaxing","poller":"twilio"}
```

## Learn More
The following can also be found at <https://godoc.org/github.com/voicera/gooseberry>
