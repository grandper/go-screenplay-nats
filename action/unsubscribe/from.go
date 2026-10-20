package unsubscribe

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// From unsubscribes the actor from the subject it subscribed to. It fails
// with action.ErrSubjectIsEmpty when the subject is empty, and with an error
// wrapping ability.ErrNoSubscriptionFound when the actor did not subscribe to
// the subject.
func From(subject string) screenplay.Performable {
	return &fromAction{subject: subject}
}

// fromAction unsubscribes the actor from a subject.
type fromAction struct {
	subject string
}

// PerformAs performs the action as the provided actor.
func (a *fromAction) PerformAs(theActor *screenplay.Actor) error {
	if a.subject == "" {
		return action.ErrSubjectIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToUnsubscribeFrom(a.subject)
}

// String describes the action.
func (a *fromAction) String() string {
	return "unsubscribe from " + verdict.Describe(a.subject)
}

// fromAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*fromAction)(nil)
