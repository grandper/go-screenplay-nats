package publish

import impl "github.com/grandper/go-screenplay-nats/v2/action/internal/publish"

// TheMessage publishes a message to the subject given with To, with the
// payload given with Containing, the headers given with WithTheHeader and
// WithTheHeaders, and the reply subject given with ToBeRepliedTo:
//
//	publish.TheMessage().To("orders").Containing([]byte("Hello")).WithTheHeader("Key", "value")
func TheMessage() *impl.TheMessageAction {
	return impl.TheMessage()
}
