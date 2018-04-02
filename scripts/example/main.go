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
