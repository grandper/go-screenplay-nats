package action

import (
	"context"

	"github.com/grandper/go-screenplay/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
)

// ReceiveAction is an action that receives a message from a NATS subject.
type ReceiveAction struct {
	subject string
}

// ReceiveFrom receives a message from a NATS subject.
func ReceiveFrom(subject string) *ReceiveAction {
	return &ReceiveAction{
		subject: subject,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (a *ReceiveAction) PerformAs(theActor *screenplay.Actor) error {
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	return communicateWithNATS.ReceiveNextMessageFor(context.Background(), a.subject)
}

// String returns a string representation of the action.
func (a *ReceiveAction) String() string {
	return "receive a message from " + a.subject
}

// ReceiveAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ReceiveAction)(nil)
