package question

import (
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// NATSMessages asks for every message the actor received, in the order they
// were received: the messages of its subscriptions and the replies to its
// requests.
func NATSMessages() screenplay.Question[[]*nats.Msg] {
	return &natsMessagesQuestion{}
}

// NATSMessagesFrom asks for the messages the actor received from the subject:
// those of its subscription to the subject, or the replies to the requests it
// sent to the subject, in the order they were received.
func NATSMessagesFrom(subject string) screenplay.Question[[]*nats.Msg] {
	return &natsMessagesQuestion{subject: subject, fromSubject: true}
}

// natsMessagesQuestion asks for the messages the actor received.
type natsMessagesQuestion struct {
	subject     string
	fromSubject bool
}

// AnsweredBy returns the messages the actor received.
func (q *natsMessagesQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*nats.Msg, error) {
	if q.fromSubject {
		return impl.TheAbilityOf(theActor).ToRetrieveMessagesFrom(q.subject)
	}

	return impl.TheAbilityOf(theActor).ToRetrieveMessages()
}

// String describes the question.
func (q *natsMessagesQuestion) String() string {
	if q.fromSubject {
		return "NATS messages from " + verdict.Describe(q.subject)
	}

	return "NATS messages"
}

// natsMessagesQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*nats.Msg] = (*natsMessagesQuestion)(nil)
