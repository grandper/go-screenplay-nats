package action

import (
	"fmt"
	"time"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// DefaultRequestTimeout is the default timeout for a NATS request.
const DefaultRequestTimeout = 5 * time.Second

// RequestAction is an action that sends a request to a NATS subject.
type RequestAction struct {
	message *ability.NATSMessage
	timeout time.Duration
	unit    string
}

// Request sends a request to a NATS subject.
func Request(message *ability.NATSMessage) *RequestAction {
	return &RequestAction{
		message: message,
		timeout: DefaultRequestTimeout,
		unit:    "second",
	}
}

// For sets the timeout for the request.
func (a *RequestAction) For(value int) *TimeBuilder {
	return &TimeBuilder{
		action: a,
		value:  value,
	}
}

type TimeBuilder struct {
	action *RequestAction
	value  int
}

// Seconds sets the timeout in seconds.
func (b *TimeBuilder) Seconds() *RequestAction {
	b.action.timeout = time.Duration(b.value) * time.Second
	b.action.unit = "second"
	return b.action
}

// Second sets the timeout in seconds.
func (b *TimeBuilder) Second() *RequestAction {
	return b.Seconds()
}

// Milliseconds sets the timeout in milliseconds.
func (b *TimeBuilder) Milliseconds() *RequestAction {
	b.action.timeout = time.Duration(b.value) * time.Millisecond
	b.action.unit = "millisecond"
	return b.action
}

// Millisecond sets the timeout in milliseconds.
func (b *TimeBuilder) Millisecond() *RequestAction {
	return b.Milliseconds()
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
	return "send a request to " + a.message.Subject() + " for " + a.timeoutAsString()
}

func (a *RequestAction) timeoutAsString() string {
	if a.unit == "second" {
		return pluralOrSingular(int(a.timeout.Seconds()), "second")
	}
	return pluralOrSingular(int(a.timeout.Milliseconds()), "millisecond")
}

func pluralOrSingular(value int, unit string) string {
	if value != 1 {
		return fmt.Sprintf("%d %ss", value, unit)
	}
	return fmt.Sprintf("%d %s", value, unit)
}

// RequestAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*RequestAction)(nil)
