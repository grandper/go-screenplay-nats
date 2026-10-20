package action

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// UnsubscribeFromAction is an action to unsubscribe from a NATS subject.
type UnsubscribeFromAction struct {
	subject string
}

// UnsubscribeFrom creates an action to unsubscribe from a NATS subject.
func UnsubscribeFrom(subject string) *UnsubscribeFromAction {
	return &UnsubscribeFromAction{
		subject: subject,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (a *UnsubscribeFromAction) PerformAs(theActor *screenplay.Actor) error {
	if a.subject == "" {
		return ErrSubjectIsEmpty
	}
	return impl.TheAbilityOf(theActor).ToUnsubscribeFrom(a.subject)
}

// String returns a string representation of the action.
func (a *UnsubscribeFromAction) String() string {
	return "unsubscribe from NATS subject " + a.subject
}

// UnsubscribeFromAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*UnsubscribeFromAction)(nil)
