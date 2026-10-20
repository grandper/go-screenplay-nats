package been_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/stretchr/testify/assert"

	questionfixture "github.com/grandper/go-screenplay-nats/v2/question/fixture"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestReconnected(t *testing.T) {
	t.Run("matches a connection that has been reconnected", func(t *testing.T) {
		fixture.AssertMatch(t, been.Reconnected(), questionfixture.CreateFakeReconnectedConnection())
	})

	t.Run("does not match a connection that has not been reconnected", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Reconnected(), questionfixture.CreateFakeConnection())
	})

	t.Run("does not match a nil connection", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Reconnected(), nil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been reconnected", been.Reconnected().String())
	})
}
