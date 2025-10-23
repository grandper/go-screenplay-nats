package question

import (
	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// NATSMessages asks about the NATS message.
func NATSMessages() *NATSMessagesQuestion {
	return &NATSMessagesQuestion{}
}

// NATSMessagesQuestion asks about the body of the NATS message.
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
func (q *NATSMessagesQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorCommunicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	if q.subject == "" {
		return theActorCommunicateWithNATS.AllMessages()
	}
	return theActorCommunicateWithNATS.AllMessagesForSubject(q.subject)
}

// NATSMessagesQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &NATSMessagesQuestion{}
