package send

import impl "github.com/grandper/go-screenplay-nats/v2/action/internal/send"

// TheRequest sends a request to the subject given with To, with the payload
// given with Containing and the headers given with WithTheHeader and
// WithTheHeaders, and waits for its reply, at most five seconds unless
// TimingOutAfter is written. The reply is stored under the subject of the
// request, where the question NATSMessagesFrom finds it:
//
//	send.TheRequest().To("orders").Containing([]byte("ping")).TimingOutAfter(2).Seconds()
func TheRequest() *impl.TheRequestAction {
	return impl.TheRequest()
}
