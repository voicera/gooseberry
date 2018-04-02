/*
Package polling provides cost-effective ways to reduce cost of polling resources.
For example: instead of busy-waiting on a service endpoint that may or may not
have a payload to process, wrapping said call as a Receiver.Receive() call and
using a PollingChannel to wrap said receiver is provides a better poller that
may use contingent on Bernoulli trials cyclic exponential backoff between calls.

For example,
    polling.NewPoller(receiver, 0.95, time.Millisecond, time.Minute)
creates a poller that keeps sending along (via the poller's send-only channel)
payloads received from the specified receiver as long as they keep arriving.
When the receiver is empty-handed (denoted by a flag on the Receive method's
return value), the pollers runs a Bernoulli trial with a probability of 0.95
to determine whether to back off (95% of the time), or call Receive() again
immediately (a costly busy-waiting; some IaaS providers, like AWS, may charge by
the call). When backing off, the poller waits between empty-handed calls; it
starts with a 1ms delay (seed), then 2ms, then 4ms, etc. The exponential delay
is capped at 1 minute (cap), after which it starts from the seed again.

The contingent nature of the cyclic exponential backoff between calls (thanks to
Bernoulli sampling) allows the poller to break out of the cycle to check for new
payloads – intermittently, which is useful when the delay intervals are long.
*/
package polling
