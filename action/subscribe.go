package action

import (
	"fmt"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// SubscribeToAction is an action to subscribe to aNATS subject.
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

// Message returns to the SubscribeToAction being built.
func (b *SubscribeToActionBuilder) Message() *SubscribeToAction {
	return b.action
}

// Messages returns to the SubscribeToAction being built.
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
func (l *SubscribeToAction) UsingTheQueue(queue string) *SubscribeToAction {
	l.queue = queue
	return l
}

// Synchronously sets the subscription to be synchronous.
func (l *SubscribeToAction) Synchronously() *SubscribeToAction {
	l.synchronous = true
	return l
}

// Asynchronously sets the subscription to be asynchronous.
func (l *SubscribeToAction) Asynchronously() *SubscribeToAction {
	l.synchronous = false
	return l
}

// AndUnsubscribeAfter sets the maximum number of wanted messages before unsubscribing.
func (l *SubscribeToAction) AndUnsubscribeAfter(maxWantedMessages int) *SubscribeToActionBuilder {
	l.maxWantedMessages = maxWantedMessages
	return &SubscribeToActionBuilder{
		action: l,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (l *SubscribeToAction) PerformAs(theActor *screenplay.Actor) error {
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	if l.subject == "" {
		return ability.ErrSubjectIsEmpty
	}
	return communicateWithNATS.SubscribeTo(l.subject, l.maxWantedMessages)
}

// String returns a string representation of the action.
func (l *SubscribeToAction) String() string {
	msg := "subscribe to NATS subject " + l.subject
	if l.synchronous {
		msg += " synchronously"
	} else {
		msg += " asynchronously"
	}
	if l.queue != "" {
		msg += " using queue " + l.queue
	}
	if l.maxWantedMessages > 0 {
		plural := "s"
		if l.maxWantedMessages == 1 {
			plural = ""
		}
		msg += fmt.Sprintf(" and unsubscribe after %d message%s", l.maxWantedMessages, plural)
	}
	return msg
}

// SubscribeToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*SubscribeToAction)(nil)
