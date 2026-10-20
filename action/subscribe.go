package action

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// SubscribeToAction is an action to subscribe to a NATS subject.
type SubscribeToAction struct {
	subject           string
	queue             string
	maxWantedMessages int
	synchronous       bool
}

// SubscribeToActionBuilder is a builder for SubscribeToAction.
type SubscribeToActionBuilder struct {
	action *SubscribeToAction
}

// Message returns the SubscribeToAction being built.
func (b *SubscribeToActionBuilder) Message() *SubscribeToAction {
	return b.action
}

// Messages returns the SubscribeToAction being built.
func (b *SubscribeToActionBuilder) Messages() *SubscribeToAction {
	return b.action
}

// SubscribeTo creates an action to subscribe to a subject on NATS.
func SubscribeTo(subject string) *SubscribeToAction {
	return &SubscribeToAction{
		subject:     subject,
		synchronous: true,
	}
}

// SubscribeToTheSubjectStartingWith creates an action to subscribe to a subject starting with a given prefix on NATS.
func SubscribeToTheSubjectStartingWith(subjectPrefix string) *SubscribeToAction {
	return &SubscribeToAction{
		subject:     subjectPrefix + ".>",
		synchronous: true,
	}
}

// UsingTheQueue sets the queue group for the subscription.
func (a *SubscribeToAction) UsingTheQueue(queue string) *SubscribeToAction {
	a.queue = queue
	return a
}

// Synchronously sets the subscription to be synchronous.
func (a *SubscribeToAction) Synchronously() *SubscribeToAction {
	a.synchronous = true
	return a
}

// Asynchronously sets the subscription to be asynchronous.
func (a *SubscribeToAction) Asynchronously() *SubscribeToAction {
	a.synchronous = false
	return a
}

// AndUnsubscribeAfter sets the maximum number of wanted messages before unsubscribing.
func (a *SubscribeToAction) AndUnsubscribeAfter(maxWantedMessages int) *SubscribeToActionBuilder {
	a.maxWantedMessages = maxWantedMessages
	return &SubscribeToActionBuilder{
		action: a,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (a *SubscribeToAction) PerformAs(theActor *screenplay.Actor) error {
	if a.subject == "" {
		return ErrSubjectIsEmpty
	}
	return impl.TheAbilityOf(theActor).ToSubscribe(impl.Subscription{
		Subject:      a.subject,
		Queue:        a.queue,
		Asynchronous: !a.synchronous,
		MaxMessages:  a.maxWantedMessages,
	})
}

// String returns a string representation of the action.
func (a *SubscribeToAction) String() string {
	msg := "subscribe to NATS subject " + a.subject
	if a.synchronous {
		msg += " synchronously"
	} else {
		msg += " asynchronously"
	}
	if a.queue != "" {
		msg += " using queue " + a.queue
	}
	if a.maxWantedMessages > 0 {
		plural := "s"
		if a.maxWantedMessages == 1 {
			plural = ""
		}
		msg += fmt.Sprintf(" and unsubscribe after %d message%s", a.maxWantedMessages, plural)
	}
	return msg
}

// SubscribeToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*SubscribeToAction)(nil)
