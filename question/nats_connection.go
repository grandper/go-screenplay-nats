package question

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// NATSConnection asks about the connection to NATS.
func NATSConnection() *NATSConnectionQuestion {
	return &NATSConnectionQuestion{}
}

// NATSConnectionQuestion asks about the NATS connection.
type NATSConnectionQuestion struct{}

// String describes the question.
func (q *NATSConnectionQuestion) String() string {
	return "connection to NATS"
}

// AnsweredBy returns a point-in-time snapshot of the connection state.
// Each call captures the state at that instant; the returned *Connection
// is immutable and becomes stale as soon as the underlying connection changes.
// Always re-ask this question rather than caching the result.
func (q *NATSConnectionQuestion) AnsweredBy(theActor *screenplay.Actor) (*ability.Connection, error) {
	return impl.TheAbilityOf(theActor).ToRetrieveTheConnection()
}

// NATSConnectionQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*ability.Connection] = (*NATSConnectionQuestion)(nil)
