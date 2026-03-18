package question

import (
	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// NATSConnection asks about the connection to NATS.
func NATSConnection() screenplay.Question {
	return &TheNATSConnectionQuestion{}
}

// TheNATSConnectionQuestion asks about the Nats connection.
type TheNATSConnectionQuestion struct{}

// String describes the question.
func (q *TheNATSConnectionQuestion) String() string {
	return "connection to NATS"
}

// AnsweredBy returns a point-in-time snapshot of the connection state.
// Each call captures the state at that instant; the returned *Connection
// is immutable and becomes stale as soon as the underlying connection changes.
// Always re-ask this question rather than caching the result.
func (q *TheNATSConnectionQuestion) AnsweredBy(theActor *screenplay.Actor) (any, error) {
	theActorCommunicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return nil, err
	}
	return &Connection{
		wasClosed:       theActorCommunicateWithNATS.ConnectionWasClosed(),
		wasDisconnected: theActorCommunicateWithNATS.ConnectionWasDisconnected(),
		wasReconnected:  theActorCommunicateWithNATS.ConnectionWasReconnected(),
	}, nil
}

// TheNATSConnectionQuestion implements the screenplay.Question interface.
var _ screenplay.Question = &TheNATSConnectionQuestion{}
