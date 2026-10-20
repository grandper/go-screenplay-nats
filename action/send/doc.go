// Package send holds the actions that send a request and wait for its reply:
//
//	err := anActor.AttemptsTo(send.TheRequest("ping").To("orders"))
//	err := anActor.AttemptsTo(send.TheBinaryRequest(data).To("orders"))
//	err := anActor.AttemptsTo(send.AnEmptyRequest().To("orders").TimingOutAfter(2).Seconds())
//
// A request carries headers with WithTheHeader and WithTheHeaders.
package send
