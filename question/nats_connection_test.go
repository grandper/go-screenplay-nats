package question_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
)

func TestNATSConnection(t *testing.T) {
	t.Run("answers with what happened to the connection", func(t *testing.T) {
		t.Run("nothing yet", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))

			connection, err := adam.AsksFor(question.NATSConnection())

			require.NoError(t, err)
			assert.False(t, connection.WasClosed())
			assert.False(t, connection.WasDisconnected())
			assert.False(t, connection.WasReconnected())
		})

		t.Run("the connection was closed", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
			nc.Close()

			assert.Eventually(t, func() bool {
				connection, err := adam.AsksFor(question.NATSConnection())
				return err == nil && connection.WasClosed() && connection.WasDisconnected() &&
					!connection.WasReconnected()
			}, timeout, tickInterval)
		})

		t.Run("the connection was reconnected", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
			require.NoError(t, nc.ForceReconnect())

			assert.Eventually(t, func() bool {
				connection, err := adam.AsksFor(question.NATSConnection())
				return err == nil && !connection.WasClosed() && connection.WasDisconnected() &&
					connection.WasReconnected()
			}, timeout, tickInterval)
		})
	})

	t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(question.NATSConnection())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("fails when the connection is not set", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
		_, err := adam.AsksFor(question.NATSConnection())
		require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "NATS connection", question.NATSConnection().String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[*ability.Connection])(nil), question.NATSConnection())
	})
}
