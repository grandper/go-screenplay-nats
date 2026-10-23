package communicate_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/ability/communicate"
	"github.com/grandper/go-screenplay-nats/v2/action/flush"
	"github.com/grandper/go-screenplay-nats/v2/action/subscribe"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
)

func TestWithNATS(t *testing.T) {
	t.Run("gives an actor the ability to communicate with NATS using a connection", func(t *testing.T) {
		nc := fixture.ConnectToAnEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

		require.NoError(t, adam.AttemptsTo(subscribe.To("orders")))
		assert.Equal(t, 1, nc.NumSubscriptions())
		require.NoError(t, adam.Exit())
		assert.True(t, nc.IsClosed())
	})

	t.Run("gives an actor the ability to communicate with NATS using a server URL", func(t *testing.T) {
		natsServer := fixture.RunAnEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").
			WhoCan(communicate.WithNATS().UsingTheServerAt(natsServer.ClientURL()))

		require.NoError(t, adam.AttemptsTo(flush.TheNATSConnection()))
		assert.Equal(t, 1, natsServer.NumClients())
		require.NoError(t, adam.Exit())
	})

	t.Run("reports an unreachable server when the actor first uses the ability", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").
			WhoCan(communicate.WithNATS().UsingTheServerAt("nats://invalid:4222"))

		require.ErrorIs(t, adam.AttemptsTo(flush.TheNATSConnection()), ability.ErrFailedToConnect)
	})

	t.Run("reports a missing connection when the actor first uses the ability", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())

		require.ErrorIs(t, adam.AttemptsTo(flush.TheNATSConnection()), ability.ErrNATSConnectionNotSet)
	})
}
