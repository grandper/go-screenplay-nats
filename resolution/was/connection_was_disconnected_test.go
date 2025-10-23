package was_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay-nats/fixture"
	"github.com/grandper/go-screenplay-nats/question"
	"github.com/grandper/go-screenplay-nats/resolution/was"
	"github.com/grandper/go-screenplay/resolution/testdata"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
)

func TestDisconnectedResolution(t *testing.T) {
	t.Run("should match if the connection is disconnected", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
		assert.NoError(t, nc.ForceReconnect())
		time.Sleep(100 * time.Millisecond)
		adam := screenplay.ActorNamed("adam").WhoCan(communicateWithNATS)
		connection, err := question.NATSConnection().AnsweredBy(adam)
		assert.NoError(t, err)

		testdata.AssertMatch(t, was.Disconnected(), connection)
	})

	t.Run("fails when the connection was not disconnected", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
		adam := screenplay.ActorNamed("adam").WhoCan(communicateWithNATS)
		connection, err := question.NATSConnection().AnsweredBy(adam)
		assert.NoError(t, err)

		testdata.AssertNoMatch(t, was.Disconnected(), connection)
		testdata.AssertNoMatch(t, was.Disconnected(), nil)
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, was.Disconnected(), 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "was disconnected", was.Disconnected().String())
	})
}
