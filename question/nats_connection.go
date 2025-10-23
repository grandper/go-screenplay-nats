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

// AnsweredBy returns the answer that an actor provided to the question.
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
