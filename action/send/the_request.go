package send

import impl "github.com/grandper/go-screenplay-nats/v2/action/internal/send"

// TheRequest sends a request holding the text to the subject given with To
// and waits for its reply, at most five seconds unless TimingOutAfter is
// written; it reads "send the request ping to orders timing out after 2
// seconds". The reply is stored under the subject of the request, where the
// question NATSMessagesFrom finds it:
//
//	send.TheRequest("ping").To("orders")
//	send.TheRequest("ping").To("orders").WithTheHeader("Key", "value").TimingOutAfter(2).Seconds()
func TheRequest(text string) *impl.TheRequestAction {
	return impl.TheRequest(text)
}

// TheBinaryRequest sends a request holding the bytes to the subject given
// with To and waits for its reply; it reads "send the binary request data to
// orders":
//
//	send.TheBinaryRequest(data).To("orders")
func TheBinaryRequest(data []byte) *impl.TheRequestAction {
	return impl.TheBinaryRequest(data)
}

// AnEmptyRequest sends a request with no payload to the subject given with
// To and waits for its reply; it reads "send an empty request to orders":
//
//	send.AnEmptyRequest().To("orders")
func AnEmptyRequest() *impl.TheRequestAction {
	return impl.AnEmptyRequest()
}
