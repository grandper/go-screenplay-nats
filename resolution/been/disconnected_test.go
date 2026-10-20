package been_test

import (
	"testing"
	"time"

	resolutionfixture "github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestDisconnectedResolution(t *testing.T) {
	t.Run("should match if the connection is disconnected", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
		require.NoError(t, nc.ForceReconnect())
		time.Sleep(100 * time.Millisecond)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS)
		connection, err := question.NATSConnection().AnsweredBy(adam)
		require.NoError(t, err)

		resolutionfixture.AssertMatch(t, been.Disconnected(), connection)
	})

	t.Run("fails when the connection was not disconnected", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS)
		connection, err := question.NATSConnection().AnsweredBy(adam)
		require.NoError(t, err)

		resolutionfixture.AssertNoMatch(t, been.Disconnected(), connection)
		resolutionfixture.AssertNoMatch(t, been.Disconnected(), nil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been disconnected", been.Disconnected().String())
	})
}
