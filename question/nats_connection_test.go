package question_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
)

func TestNATSConnection(t *testing.T) {
	t.Run("returns the NATS connection data", func(t *testing.T) {
		t.Run("no events have occurred", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			connection, err := question.NATSConnection().AnsweredBy(adam)
			require.NoError(t, err)
			assert.False(t, connection.WasClosed())
			assert.False(t, connection.WasDisconnected())
			assert.False(t, connection.WasReconnected())
		})

		t.Run("the connection was closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			nc.Close()
			time.Sleep(100 * time.Millisecond)
			connection, err := question.NATSConnection().AnsweredBy(adam)
			require.NoError(t, err)
			assert.True(t, connection.WasClosed())
			assert.True(t, connection.WasDisconnected())
			assert.False(t, connection.WasReconnected())
		})

		t.Run("the connection was reconnected", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.NoError(t, nc.ForceReconnect())
			time.Sleep(100 * time.Millisecond)
			connection, err := question.NATSConnection().AnsweredBy(adam)
			require.NoError(t, err)
			assert.False(t, connection.WasClosed())
			assert.True(t, connection.WasDisconnected())
			assert.True(t, connection.WasReconnected())
		})
	})

	t.Run(
		"fails to get the NATS connection data if the actor does not have the ability to communicate with NATS",
		func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			_, err := question.NATSConnection().AnsweredBy(adam)
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		},
	)

	t.Run("implements the stringer interface", func(t *testing.T) {
		body := question.NATSConnection()
		assert.Equal(t, "connection to NATS", body.String())
	})
}
