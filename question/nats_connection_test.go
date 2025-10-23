package question_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay-nats/fixture"
	"github.com/grandper/go-screenplay-nats/question"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNATSConnection(t *testing.T) {
	t.Run("returns the NATS connection data", func(t *testing.T) {
		t.Run("no events have occurred", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			connection, err := question.NATSConnection().AnsweredBy(adam)
			require.NoError(t, err)
			assert.False(t, connection.(*question.Connection).WasClosed())
			assert.False(t, connection.(*question.Connection).WasDisconnected())
			assert.False(t, connection.(*question.Connection).WasReconnected())
		})

		t.Run("the connection was closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			nc.Close()
			time.Sleep(100 * time.Millisecond)
			connection, err := question.NATSConnection().AnsweredBy(adam)
			require.NoError(t, err)
			assert.True(t, connection.(*question.Connection).WasClosed())
			assert.True(t, connection.(*question.Connection).WasDisconnected())
			assert.False(t, connection.(*question.Connection).WasReconnected())
		})

		t.Run("the connection was reconnected", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.NoError(t, nc.ForceReconnect())
			time.Sleep(100 * time.Millisecond)
			connection, err := question.NATSConnection().AnsweredBy(adam)
			require.NoError(t, err)
			assert.False(t, connection.(*question.Connection).WasClosed())
			assert.True(t, connection.(*question.Connection).WasDisconnected())
			assert.True(t, connection.(*question.Connection).WasReconnected())
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
