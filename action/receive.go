package action

import (
	"context"

	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// ReceiveAction is an action that receives a message from a NATS subject.
type ReceiveAction struct {
	subject string
}

// ReceiveFrom receives the next message from a subject the actor subscribed
// to synchronously, waiting at most the timeout of the actor.
func ReceiveFrom(subject string) *ReceiveAction {
	return &ReceiveAction{
		subject: subject,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (a *ReceiveAction) PerformAs(theActor *screenplay.Actor) error {
	if a.subject == "" {
		return ErrSubjectIsEmpty
	}

	ctx, cancel := context.WithTimeout(theActor.Context(), theActor.Timeout())
	defer cancel()

	return impl.TheAbilityOf(theActor).ToReceiveTheNextMessageFrom(ctx, a.subject)
}

// String returns a string representation of the action.
func (a *ReceiveAction) String() string {
	return "receive a message from " + a.subject
}

// ReceiveAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ReceiveAction)(nil)
