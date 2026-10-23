package receive

import (
	"context"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// TheNextMessageFrom receives the next message of the synchronous subscription
// to the subject and stores it, where messages.FromTheSubject finds it.
// The actor waits at most its timeout (see screenplay.WithTimeout). It fails
// with action.ErrSubjectIsEmpty when the subject is empty, and with an error
// wrapping ability.ErrNoSubscriptionFound when the actor did not subscribe to
// the subject.
func TheNextMessageFrom(subject string) screenplay.Performable {
	return &theNextMessageFromAction{subject: subject}
}

// theNextMessageFromAction receives the next message of a subscription.
type theNextMessageFromAction struct {
	subject string
}

// PerformAs performs the action as the provided actor.
func (a *theNextMessageFromAction) PerformAs(theActor *screenplay.Actor) error {
	if a.subject == "" {
		return action.ErrSubjectIsEmpty
	}

	ctx, cancel := context.WithTimeout(theActor.Context(), theActor.Timeout())
	defer cancel()

	return impl.TheAbilityOf(theActor).ToReceiveTheNextMessageFrom(ctx, a.subject)
}

// String describes the action.
func (a *theNextMessageFromAction) String() string {
	return "receive the next message from " + verdict.Describe(a.subject)
}

// theNextMessageFromAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*theNextMessageFromAction)(nil)
