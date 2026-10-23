package fixture

import (
	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/nats-io/nats.go"

	"github.com/grandper/go-screenplay-nats/v2/ability/communicate"
)

// NarratedActor creates an actor whose narration is recorded by the recorder,
// and who can communicate with NATS using the connection.
func NarratedActor(name string, recorder *screenplayfixture.Recorder, nc *nats.Conn) *screenplay.Actor {
	production := screenplay.NewProduction(screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

	return production.SetTheStage(screenplay.CastOfStandardActors()).
		ActorNamed(name).
		WhoCan(communicate.WithNATS().UsingTheConnection(nc))
}
