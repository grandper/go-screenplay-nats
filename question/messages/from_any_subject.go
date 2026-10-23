package messages

import (
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// FromAnySubject asks for every message the actor received, in the order they
// were received: the messages of its subscriptions and the replies to its
// requests.
func FromAnySubject() screenplay.Question[[]*nats.Msg] {
	return &fromAnySubjectQuestion{}
}

// fromAnySubjectQuestion asks for every message the actor received.
type fromAnySubjectQuestion struct{}

// AnsweredBy returns every message the actor received.
func (q *fromAnySubjectQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*nats.Msg, error) {
	return impl.TheAbilityOf(theActor).ToRetrieveMessages()
}

// String describes the question.
func (q *fromAnySubjectQuestion) String() string {
	return "NATS messages from any subject"
}

// fromAnySubjectQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*nats.Msg] = (*fromAnySubjectQuestion)(nil)
