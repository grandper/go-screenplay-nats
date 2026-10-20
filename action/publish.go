package action

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// PublishAction is an action that publishes a message to a NATS subject.
type PublishAction struct {
	message *ability.NATSMessage
}

// Publish publishes a message to a NATS subject.
func Publish(message *ability.NATSMessage) *PublishAction {
	return &PublishAction{
		message: message,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (a *PublishAction) PerformAs(theActor *screenplay.Actor) error {
	if a.message != nil && a.message.Subject() == "" {
		return ErrSubjectIsEmpty
	}
	return impl.TheAbilityOf(theActor).ToPublish(a.message)
}

// String returns a string representation of the action.
func (a *PublishAction) String() string {
	if a.message == nil {
		return "publish a nil message"
	}
	return "publish a message to " + a.message.Subject()
}

// PublishAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*PublishAction)(nil)
