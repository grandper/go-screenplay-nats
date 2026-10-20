package subscribe

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// To subscribes the actor to the subject, synchronously unless Asynchronously
// is written: the messages of a synchronous subscription are received with
// receive.TheNextMessageFrom, those of an asynchronous one are stored as
// they arrive. Wildcards are allowed in the subject:
//
//	subscribe.To("orders.*.created")
//	subscribe.To("orders").InTheQueueGroup("workers").Asynchronously()
//	subscribe.To("orders").AndUnsubscribeAfter(10).Messages()
func To(subject string) *ToAction {
	return &ToAction{
		subject:     subject,
		description: "subscribe to " + verdict.Describe(subject),
	}
}

// ToTheSubjectsStartingWith subscribes the actor to every subject starting
// with the prefix, which is the subject "prefix.>"; it reads "subscribe to
// the subjects starting with orders":
//
//	subscribe.ToTheSubjectsStartingWith("orders")
func ToTheSubjectsStartingWith(prefix string) *ToAction {
	subject := prefix + ".>"
	if prefix == "" {
		subject = ""
	}

	return &ToAction{
		subject:     subject,
		description: "subscribe to the subjects starting with " + verdict.Describe(prefix),
	}
}

// ToAction subscribes the actor to a subject.
type ToAction struct {
	subject      string
	queue        string
	queueWritten bool
	asynchronous bool
	maxMessages  int
	description  string
	parts        []string
}

// InTheQueueGroup makes the subscription part of the queue group, so that
// each message is delivered to one subscriber of the group; it reads "in the
// queue group workers".
func (a *ToAction) InTheQueueGroup(queue string) *ToAction {
	a.queue = queue
	a.queueWritten = true
	a.parts = append(a.parts, "in the queue group "+verdict.Describe(queue))

	return a
}

// Synchronously makes the subscription synchronous: its messages are received
// with receive.TheNextMessageFrom. It is the default.
func (a *ToAction) Synchronously() *ToAction {
	a.asynchronous = false
	a.parts = append(a.parts, "synchronously")

	return a
}

// Asynchronously makes the subscription asynchronous: its messages are stored
// as they arrive.
func (a *ToAction) Asynchronously() *ToAction {
	a.asynchronous = true
	a.parts = append(a.parts, "asynchronously")

	return a
}

// AndUnsubscribeAfter bounds the subscription to a number of messages; it reads
// "and unsubscribe after 10 messages".
func (a *ToAction) AndUnsubscribeAfter(amount int) *UnsubscribeAfterBuilder {
	return &UnsubscribeAfterBuilder{action: a, amount: amount}
}

// UnsubscribeAfterBuilder holds the amount given to AndUnsubscribeAfter and
// waits for the unit that closes the setting.
type UnsubscribeAfterBuilder struct {
	action *ToAction
	amount int
}

// Messages makes the amount the number of messages after which the actor
// unsubscribes.
func (b *UnsubscribeAfterBuilder) Messages() *ToAction {
	b.action.maxMessages = b.amount
	unit := "messages"
	if b.amount == 1 {
		unit = "message"
	}
	b.action.parts = append(b.action.parts, fmt.Sprintf("and unsubscribe after %d %s", b.amount, unit))

	return b.action
}

// Message is the singular of Messages, so that one message reads as English.
func (b *UnsubscribeAfterBuilder) Message() *ToAction {
	return b.Messages()
}

// PerformAs performs the action as the provided actor.
func (a *ToAction) PerformAs(theActor *screenplay.Actor) error {
	if a.subject == "" {
		return action.ErrSubjectIsEmpty
	}
	if a.queueWritten && a.queue == "" {
		return action.ErrQueueIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToSubscribeTo(a.subject, impl.SubscriptionOptions{
		Queue:        a.queue,
		Asynchronous: a.asynchronous,
		MaxMessages:  a.maxMessages,
	})
}

// String describes the action, word for word as it was written.
func (a *ToAction) String() string {
	return strings.Join(append([]string{a.description}, a.parts...), " ")
}

// ToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ToAction)(nil)
