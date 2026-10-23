package messages

import (
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// FromTheSubject asks for the messages the actor received from the subject:
// those of its subscription to the subject, or the replies to the requests it
// sent to the subject, in the order they were received.
func FromTheSubject(subject string) screenplay.Question[[]*nats.Msg] {
	return &fromTheSubjectQuestion{subject: subject}
}

// fromTheSubjectQuestion asks for the messages the actor received from a subject.
type fromTheSubjectQuestion struct {
	subject string
}

// AnsweredBy returns the messages the actor received from the subject.
func (q *fromTheSubjectQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*nats.Msg, error) {
	return impl.TheAbilityOf(theActor).ToRetrieveMessagesFrom(q.subject)
}

// String describes the question.
func (q *fromTheSubjectQuestion) String() string {
	return "NATS messages from the subject " + verdict.Describe(q.subject)
}

// fromTheSubjectQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*nats.Msg] = (*fromTheSubjectQuestion)(nil)
