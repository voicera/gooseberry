package polling

import (
	"time"

	"github.com/voicera/gooseberry"
)

// Poller represents the resource being polled as a send-only channel.
// Ideally, this would be a generic type (see https://golang.org/doc/faq#generics).
type Poller interface {
	// Channel gets the underlying channel that wraps the resource being polled.
	// By default, sends block until the other side (the receiver) is ready.
	Channel() <-chan interface{}

	// GetName returns name of the poller.
	GetName() string

	// Start starts polling for new payloads to arrive on the receiving end.
	// This method is non-idempotent.
	Start()

	// Stop signals the poller to stop polling; it also closes the channel.
	// This method is non-idempotent.
	Stop()
}

// Receiver represents the simple act of receiving payloads (e.g., messages)
// from the resource being polled.
type Receiver interface {
	// Receive receives from the resource being polled, returning
	// the received payload (if any), a flag to indicate whether or not
	// a payload was found, and an error if encountered.
	Receive() (interface{}, bool, error)
}

type pollingChannel struct {
	name     string
	receiver Receiver
	data     chan interface{}
	signal   chan interface{} // can be user later to send commands (like relax)
	relaxer  relaxer
}

// NewBernoulliExponentialBackoffPoller creates a polling channel that uses
// Bernoulli trials with cyclic exponential backoff between empty-handed calls.
func NewBernoulliExponentialBackoffPoller(
	receiver Receiver, entityName string, probability float64, seed, cap time.Duration) (Poller, error) {
	sampler, err := newBernoulliSampler(probability)
	if err != nil {
		return nil, err
	}
	return &pollingChannel{
		name:     entityName,
		data:     make(chan interface{}),
		signal:   make(chan interface{}),
		receiver: receiver,
		relaxer: &cyclicExponentialBackoffRelaxer{
			relaxationCondition:    &and{[]relaxationCondition{&emptyHanded{}, sampler}},
			initialBackoffDuration: seed,
			currentBackoffDuration: seed,
			exponentialBackoffCap:  cap,
			sleep: time.Sleep,
		},
	}, nil
}

func (pc *pollingChannel) Channel() <-chan interface{} {
	return pc.data
}

func (pc *pollingChannel) GetName() string {
	return pc.name
}

func (pc *pollingChannel) Start() {
	gooseberry.Logger.Debug("Started", "poller", pc.name)
	defer gooseberry.Logger.Debug("Stopped", "poller", pc.name)
	defer close(pc.data)

	for {
		if payload, found, err := pc.receiver.Receive(); found && err == nil {
			select {
			case pc.data <- payload:
				pc.relax(found)
			case <-pc.signal:
				return
			}
		} else {
			if err != nil {
				gooseberry.Logger.Error(err.Error(), "poller", pc.name, "err", err)
			}

			select {
			case <-pc.signal:
				return
			default:
				pc.relax(false)
			}
		}
	}
}

func (pc *pollingChannel) Stop() {
	gooseberry.Logger.Debug("Stopping", "poller", pc.name)
	pc.signal <- 0
}

func (pc *pollingChannel) relax(lastReceiptSucceeded bool) {
	gooseberry.Logger.Debug("Relaxing", "poller", pc.name)
	pc.relaxer.relax(lastReceiptSucceeded)
}
