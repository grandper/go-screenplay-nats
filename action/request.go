package action

import (
	"time"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/grandper/go-screenplay/timing"
)

// DefaultRequestTimeout is the default timeout for a NATS request.
const DefaultRequestTimeout = 5 * time.Second

// RequestAction is an action that sends a request to a NATS subject.
type RequestAction struct {
	durationBuilder *timing.DurationBuilder[RequestAction]
	message         *ability.NATSMessage
	timeout         time.Duration
}

// Request sends a request to a NATS subject.
func Request(message *ability.NATSMessage) *RequestAction {
	action := &RequestAction{
		message: message,
		timeout: DefaultRequestTimeout,
	}
	action.durationBuilder = timing.NewDurationBuilder(action, &action.timeout)
	action.durationBuilder.For(int(DefaultRequestTimeout / time.Second)).Seconds()

	return action
}

// For sets the timeout for the request. The returned DurationBuilder lets the
// caller pick a unit, for example Request(msg).For(20).Milliseconds().
func (a *RequestAction) For(value int) *timing.DurationBuilder[RequestAction] {
	return a.durationBuilder.For(value)
}

// PerformAs performs the task or the action as the provided actor.
func (a *RequestAction) PerformAs(theActor *screenplay.Actor) error {
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	_, err = communicateWithNATS.Request(a.message, a.timeout)
	return err
}

// String returns a string representation of the action.
func (a *RequestAction) String() string {
	if a.message == nil {
		return "send a nil request"
	}
	return "send a request to " + a.message.Subject() + " for " + a.durationBuilder.String()
}

// RequestAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*RequestAction)(nil)
