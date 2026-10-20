package been_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/stretchr/testify/assert"

	questionfixture "github.com/grandper/go-screenplay-nats/v2/question/fixture"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestDisconnected(t *testing.T) {
	t.Run("matches a connection that has been disconnected", func(t *testing.T) {
		fixture.AssertMatch(t, been.Disconnected(), questionfixture.CreateFakeDisconnectedConnection())
	})

	t.Run(
		"matches a connection that has been closed or reconnected, since it was disconnected first",
		func(t *testing.T) {
			fixture.AssertMatch(t, been.Disconnected(), questionfixture.CreateFakeClosedConnection())
			fixture.AssertMatch(t, been.Disconnected(), questionfixture.CreateFakeReconnectedConnection())
		},
	)

	t.Run("does not match a connection that has not been disconnected", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Disconnected(), questionfixture.CreateFakeConnection())
	})

	t.Run("does not match a nil connection", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Disconnected(), nil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been disconnected", been.Disconnected().String())
	})
}
