// Package publish holds the actions that publish a message:
//
//	err := anActor.AttemptsTo(publish.TheMessage("Hello").To("orders"))
//	err := anActor.AttemptsTo(publish.TheBinaryMessage(data).To("orders"))
//	err := anActor.AttemptsTo(publish.AnEmptyMessage().To("orders"))
//
// A message carries headers with WithTheHeader and WithTheHeaders, and the
// subject a reply is expected on with ExpectingAReplyOn.
package publish
