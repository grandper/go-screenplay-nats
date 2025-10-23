package action

import (
	"errors"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrSubjectIsEmpty is returned when the subject for the message is empty.
	ErrSubjectIsEmpty = errors.New("subject is empty")
)

// PublishAction is an action that publishes a message to a NATS subject.
type PublishAction struct {
	message *ability.NATSMessage
}

// Publish publishes a message to a NATS.
func Publish(message *ability.NATSMessage) *PublishAction {
	return &PublishAction{
		message: message,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (s *PublishAction) PerformAs(theActor *screenplay.Actor) error {
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	return communicateWithNATS.Publish(s.message)
}

// String returns a string representation of the action.
func (s *PublishAction) String() string {
	return "publish a message to " + s.message.Subject()
}

// PublishAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*PublishAction)(nil)
