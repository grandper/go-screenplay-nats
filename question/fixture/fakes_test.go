package fixture_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay-nats/v2/question/fixture"
)

func TestFakes(t *testing.T) {
	t.Run("creates a fake connection nothing happened to", func(t *testing.T) {
		connection := fixture.CreateFakeConnection()
		assert.False(t, connection.WasClosed())
		assert.False(t, connection.WasDisconnected())
		assert.False(t, connection.WasReconnected())
	})

	t.Run("creates a fake closed connection", func(t *testing.T) {
		connection := fixture.CreateFakeClosedConnection()
		assert.True(t, connection.WasClosed())
		assert.True(t, connection.WasDisconnected())
		assert.False(t, connection.WasReconnected())
	})

	t.Run("creates a fake disconnected connection", func(t *testing.T) {
		connection := fixture.CreateFakeDisconnectedConnection()
		assert.False(t, connection.WasClosed())
		assert.True(t, connection.WasDisconnected())
		assert.False(t, connection.WasReconnected())
	})

	t.Run("creates a fake reconnected connection", func(t *testing.T) {
		connection := fixture.CreateFakeReconnectedConnection()
		assert.False(t, connection.WasClosed())
		assert.True(t, connection.WasDisconnected())
		assert.True(t, connection.WasReconnected())
	})
}
