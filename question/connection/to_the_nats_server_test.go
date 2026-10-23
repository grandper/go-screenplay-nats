package connection_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/ability/communicate"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question/connection"
)

const (
	timeout      = 2 * time.Second
	tickInterval = 100 * time.Millisecond
)

func TestToTheNATSServer(t *testing.T) {
	t.Run("answers with what happened to the connection", func(t *testing.T) {
		t.Run("nothing yet", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			answer, err := adam.AsksFor(connection.ToTheNATSServer())

			require.NoError(t, err)
			assert.False(t, answer.WasClosed())
			assert.False(t, answer.WasDisconnected())
			assert.False(t, answer.WasReconnected())
		})

		t.Run("the connection was closed", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			nc.Close()

			assert.Eventually(t, func() bool {
				answer, err := adam.AsksFor(connection.ToTheNATSServer())
				return err == nil && answer.WasClosed() && answer.WasDisconnected() && !answer.WasReconnected()
			}, timeout, tickInterval)
		})

		t.Run("the connection was reconnected", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			require.NoError(t, nc.ForceReconnect())

			assert.Eventually(t, func() bool {
				answer, err := adam.AsksFor(connection.ToTheNATSServer())
				return err == nil && !answer.WasClosed() && answer.WasDisconnected() && answer.WasReconnected()
			}, timeout, tickInterval)
		})
	})

	t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(connection.ToTheNATSServer())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("fails when the connection is not set", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
		_, err := adam.AsksFor(connection.ToTheNATSServer())
		require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "connection to the NATS server", connection.ToTheNATSServer().String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[*ability.Connection])(nil), connection.ToTheNATSServer())
	})
}
