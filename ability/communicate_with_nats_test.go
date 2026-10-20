package ability_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action/flush"
	"github.com/grandper/go-screenplay-nats/v2/action/subscribe"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
)

func TestCommunicateWithNATS(t *testing.T) {
	t.Run("gives an actor the ability to communicate with NATS using a connection", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))

		require.NoError(t, adam.AttemptsTo(subscribe.To("orders")))
		assert.Equal(t, 1, nc.NumSubscriptions())
		require.NoError(t, adam.Exit())
		assert.True(t, nc.IsClosed())
	})

	t.Run("gives an actor the ability to communicate with NATS using a server URL", func(t *testing.T) {
		natsServer := fixture.RunNATSServer(t)
		adam := screenplay.ActorNamed("Adam").
			WhoCan(ability.CommunicateWithNATS().UsingServerAt(natsServer.ClientURL()))

		require.NoError(t, adam.AttemptsTo(flush.TheNATSConnection()))
		assert.Equal(t, 1, natsServer.NumClients())
		require.NoError(t, adam.Exit())
	})

	t.Run("reports an unreachable server when the actor first uses the ability", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingServerAt("nats://invalid:4222"))

		require.ErrorIs(t, adam.AttemptsTo(flush.TheNATSConnection()), ability.ErrFailedToConnect)
	})

	t.Run("reports a missing connection when the actor first uses the ability", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())

		require.ErrorIs(t, adam.AttemptsTo(flush.TheNATSConnection()), ability.ErrNATSConnectionNotSet)
	})
}
