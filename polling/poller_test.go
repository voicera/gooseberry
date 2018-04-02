package polling

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type counter struct{ count int }
type alwaysEmptyHandedReceiver counter
type neverEmptyHandedReceiver counter
type threeSidedDieReceiver counter

func (*alwaysEmptyHandedReceiver) Receive() (payload interface{}, found bool, err error) {
	return
}

func (receiver *neverEmptyHandedReceiver) Receive() (interface{}, bool, error) {
	receiver.count++
	return receiver.count, true, nil
}

func (receiver *threeSidedDieReceiver) Receive() (interface{}, bool, error) {
	receiver.count++
	cast := rand.Float64()
	if cast < 1.0/3 {
		return receiver.count, true, nil
	} else if cast < 2.0/3 {
		return receiver.count, false, nil
	} else {
		return nil, false, errors.New("failed successfully")
	}
}

func ExamplePoller() {
	rand.Seed(0) // to produce the same sequence of pseudo-random numbers every time
	receivers := []Receiver{&alwaysEmptyHandedReceiver{}, &neverEmptyHandedReceiver{}, &threeSidedDieReceiver{}}

	for r, receiver := range receivers {
		poller, err := NewBernoulliExponentialBackoffPoller(
			receiver, "test", 0.5, time.Nanosecond, time.Millisecond)
		if err != nil {
			panic(err)
		}
		payloads := []interface{}{}

		go poller.Start()

		for i := 0; r > 0 && i < 13; i++ {
			payloads = append(payloads, <-poller.Channel())
		}

		poller.Stop()
		_, receiving := <-poller.Channel()
		fmt.Printf("Stopped; is channel receiving from %T? %t\n", receiver, receiving)
		fmt.Println("Payloads received:", payloads)
		fmt.Println()
	}
	// Output:
	// Stopped; is channel receiving from *polling.alwaysEmptyHandedReceiver? false
	// Payloads received: []
	//
	// Stopped; is channel receiving from *polling.neverEmptyHandedReceiver? false
	// Payloads received: [1 2 3 4 5 6 7 8 9 10 11 12 13]
	//
	// Stopped; is channel receiving from *polling.threeSidedDieReceiver? false
	// Payloads received: [4 6 7 9 15 17 19 20 21 25 31 38 44]
}
