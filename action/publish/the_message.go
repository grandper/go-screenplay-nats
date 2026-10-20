package publish

import impl "github.com/grandper/go-screenplay-nats/v2/action/internal/publish"

// TheMessage publishes a message holding the text to the subject given with
// To; it reads "publish the message Hello to orders":
//
//	publish.TheMessage("Hello").To("orders")
//	publish.TheMessage("Hello").To("orders").WithTheHeader("Key", "value").ExpectingAReplyOn("inbox")
func TheMessage(text string) *impl.TheMessageAction {
	return impl.TheMessage(text)
}

// TheBinaryMessage publishes a message holding the bytes to the subject given
// with To; it reads "publish the binary message data to orders":
//
//	publish.TheBinaryMessage(data).To("orders")
func TheBinaryMessage(data []byte) *impl.TheMessageAction {
	return impl.TheBinaryMessage(data)
}

// AnEmptyMessage publishes a message with no payload to the subject given
// with To; it reads "publish an empty message to orders":
//
//	publish.AnEmptyMessage().To("orders")
func AnEmptyMessage() *impl.TheMessageAction {
	return impl.AnEmptyMessage()
}
