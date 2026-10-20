package publish

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action/internal/message"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// TheMessage publishes a message holding the text to the subject given with
// To; it reads "publish the message Hello to orders":
//
//	publish.TheMessage("Hello").To("orders")
//	publish.TheMessage("Hello").To("orders").WithTheHeader("Key", "value").ExpectingAReplyOn("inbox")
func TheMessage(text string) *TheMessageAction {
	return &TheMessageAction{
		description: "publish the message " + verdict.Describe(text),
		message:     message.NewBuilder([]byte(text)),
	}
}

// TheBinaryMessage publishes a message holding the bytes to the subject given
// with To; it reads "publish the binary message data to orders":
//
//	publish.TheBinaryMessage(data).To("orders")
func TheBinaryMessage(data []byte) *TheMessageAction {
	return &TheMessageAction{
		description: "publish the binary message " + verdict.Describe(data),
		message:     message.NewBuilder(data),
	}
}

// AnEmptyMessage publishes a message with no payload to the subject given
// with To; it reads "publish an empty message to orders":
//
//	publish.AnEmptyMessage().To("orders")
func AnEmptyMessage() *TheMessageAction {
	return &TheMessageAction{
		description: "publish an empty message",
		message:     message.NewBuilder(nil),
	}
}

// TheMessageAction publishes a message.
type TheMessageAction struct {
	description string
	message     *message.Builder
}

// To sets the subject the message is published to; it reads "to orders".
func (a *TheMessageAction) To(subject string) *TheMessageAction {
	a.message.To(subject)

	return a
}

// WithTheHeader adds a header to the message; it reads "with the header
// Key = value".
func (a *TheMessageAction) WithTheHeader(key, value string) *TheMessageAction {
	a.message.WithTheHeader(key, value)

	return a
}

// WithTheHeaders adds headers to the message, written as key, value, key,
// value; it reads "with the headers A = 1, B = 2".
func (a *TheMessageAction) WithTheHeaders(args ...string) *TheMessageAction {
	a.message.WithTheHeaders(args...)

	return a
}

// ExpectingAReplyOn sets the subject the message is to be replied on; it
// reads "expecting a reply on inbox".
func (a *TheMessageAction) ExpectingAReplyOn(reply string) *TheMessageAction {
	a.message.ExpectingAReplyOn(reply)

	return a
}

// PerformAs performs the action as the provided actor. It fails with
// action.ErrSubjectIsEmpty when no subject was given, and with an error
// wrapping action.ErrInvalidHeaders when the headers were written wrongly.
func (a *TheMessageAction) PerformAs(theActor *screenplay.Actor) error {
	if err := a.message.Err(); err != nil {
		return err
	}

	return impl.TheAbilityOf(theActor).ToPublish(a.message.Msg())
}

// String describes the action, word for word as it was written.
func (a *TheMessageAction) String() string {
	if chain := a.message.Describe(); chain != "" {
		return a.description + " " + chain
	}

	return a.description
}

// TheMessageAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheMessageAction)(nil)
