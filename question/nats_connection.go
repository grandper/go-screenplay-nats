package question

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// NATSConnection asks what happened to the NATS connection of the actor. The
// answer is a snapshot taken when the question is answered; ask again to
// observe a later state.
func NATSConnection() screenplay.Question[*ability.Connection] {
	return &natsConnectionQuestion{}
}

// natsConnectionQuestion asks what happened to the NATS connection of the actor.
type natsConnectionQuestion struct{}

// AnsweredBy returns a snapshot of what happened to the connection.
func (q *natsConnectionQuestion) AnsweredBy(theActor *screenplay.Actor) (*ability.Connection, error) {
	return impl.TheAbilityOf(theActor).ToTellWhatHappenedToTheConnection()
}

// String describes the question.
func (q *natsConnectionQuestion) String() string {
	return "NATS connection"
}

// natsConnectionQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*ability.Connection] = (*natsConnectionQuestion)(nil)
