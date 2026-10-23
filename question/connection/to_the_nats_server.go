package connection

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// ToTheNATSServer asks what happened to the connection of the actor to the
// NATS server. The answer is a snapshot taken when the question is answered;
// ask again to observe a later state.
func ToTheNATSServer() screenplay.Question[*ability.Connection] {
	return &toTheNATSServerQuestion{}
}

// toTheNATSServerQuestion asks what happened to the connection of the actor to the NATS server.
type toTheNATSServerQuestion struct{}

// AnsweredBy returns a snapshot of what happened to the connection.
func (q *toTheNATSServerQuestion) AnsweredBy(theActor *screenplay.Actor) (*ability.Connection, error) {
	return impl.TheAbilityOf(theActor).ToTellWhatHappenedToTheConnection()
}

// String describes the question.
func (q *toTheNATSServerQuestion) String() string {
	return "connection to the NATS server"
}

// toTheNATSServerQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*ability.Connection] = (*toTheNATSServerQuestion)(nil)
