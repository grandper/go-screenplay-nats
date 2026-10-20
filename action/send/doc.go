// Package send holds the action that sends a request and waits for its reply:
//
//	err := anActor.AttemptsTo(send.TheRequest().To("orders").Containing([]byte("ping")))
//	err := anActor.AttemptsTo(send.TheRequest().To("orders").TimingOutAfter(2).Seconds())
package send
