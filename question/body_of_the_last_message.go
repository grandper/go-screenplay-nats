package question

import (
	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// BodyOfTheLastMessage asks about the body of the last NATS message.
func BodyOfTheLastMessage() *BodyOfTheLastMessageQuestion {
	return &BodyOfTheLastMessageQuestion{}
}

// BodyOfTheLastMessageQuestion asks about the body of the last Nats message.
type BodyOfTheLastMessageQuestion struct {
	subject string
}

// From sets the subject from which to get the last message.
func (q *BodyOfTheLastMessageQuestion) From(subject string) *BodyOfTheLastMessageQuestion {
	q.subject = subject
	return q
}

// String describes the question.
func (q *BodyOfTheLastMessageQuestion) String() string {
	if q.subject != "" {
		return "body of the last message from " + q.subject
	}
	return "body of the last message"
}

// AnsweredBy returns the answer that an actor provided to the question.
func (q *BodyOfTheLastMessageQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorCommunicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	msg, err := theActorCommunicateWithNATS.LastMessageForSubject(q.subject)
	if err != nil {
		return nil, err
	}
	return string(msg.Data), nil
}

// BodyOfTheLastMessageQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &BodyOfTheLastMessageQuestion{}
