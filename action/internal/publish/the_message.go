package publish

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/internal/message"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// TheMessage publishes a message to the subject given with To, with the
// payload given with Containing, the headers given with WithTheHeader and
// WithTheHeaders, and the reply subject given with ToBeRepliedTo:
//
//	publish.TheMessage().To("orders").Containing([]byte("Hello")).WithTheHeader("Key", "value")
func TheMessage() *TheMessageAction {
	a := &TheMessageAction{}
	a.Builder = message.NewBuilder(a)

	return a
}

// TheMessageAction publishes a message. Its chain methods To, Containing,
// WithTheHeader, WithTheHeaders and ToBeRepliedTo come from the message
// builder it embeds.
type TheMessageAction struct {
	*message.Builder[TheMessageAction]
}

// PerformAs performs the action as the provided actor. It fails with
// action.ErrSubjectIsEmpty when no subject was given.
func (a *TheMessageAction) PerformAs(theActor *screenplay.Actor) error {
	if a.Subject() == "" {
		return action.ErrSubjectIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToPublish(a.Msg())
}

// String describes the action, word for word as it was written.
func (a *TheMessageAction) String() string {
	if chain := a.Describe(); chain != "" {
		return "publish the message " + chain
	}

	return "publish the message"
}

// TheMessageAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheMessageAction)(nil)
