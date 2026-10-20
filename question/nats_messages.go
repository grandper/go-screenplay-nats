package question

import (
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// NATSMessages asks about the NATS messages received by the actor.
func NATSMessages() *NATSMessagesQuestion {
	return &NATSMessagesQuestion{}
}

// NATSMessagesQuestion asks about the NATS messages received by the actor.
type NATSMessagesQuestion struct {
	subject string
}

// From sets the subject from which to get the messages.
func (q *NATSMessagesQuestion) From(subject string) *NATSMessagesQuestion {
	q.subject = subject
	return q
}

// String describes the question.
func (q *NATSMessagesQuestion) String() string {
	if q.subject != "" {
		return "NATS messages from " + q.subject
	}
	return "NATS messages"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *NATSMessagesQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*nats.Msg, error) {
	if q.subject == "" {
		return impl.TheAbilityOf(theActor).ToRetrieveMessages()
	}
	return impl.TheAbilityOf(theActor).ToRetrieveMessagesFrom(q.subject)
}

// NATSMessagesQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*nats.Msg] = (*NATSMessagesQuestion)(nil)
