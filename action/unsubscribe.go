package action

import (
	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
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
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	if a.subject == "" {
		return ability.ErrSubjectIsEmpty
	}
	return communicateWithNATS.UnsubscribeFrom(a.subject)
}

// String returns a string representation of the action.
func (a *UnsubscribeFromAction) String() string {
	return "unsubscribe from NATS subject " + a.subject
}

// UnsubscribeFromAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*UnsubscribeFromAction)(nil)
