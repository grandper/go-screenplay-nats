package been_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	questionfixture "github.com/grandper/go-screenplay-nats/v2/question/fixture"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestConnection(t *testing.T) {
	nothingHappened := questionfixture.CreateFakeConnection()
	closed := questionfixture.CreateFakeClosedConnection()
	disconnected := questionfixture.CreateFakeDisconnectedConnection()
	reconnected := questionfixture.CreateFakeReconnectedConnection()

	t.Run("been closed", func(t *testing.T) {
		t.Run("matches a connection that has been closed", func(t *testing.T) {
			fixture.AssertMatch(t, been.Closed(), closed)
		})

		t.Run("does not match a connection that has not been closed", func(t *testing.T) {
			fixture.AssertNoMatch(t, been.Closed(), nothingHappened)
			fixture.AssertNoMatch(t, been.Closed(), disconnected)
			fixture.AssertNoMatch(t, been.Closed(), reconnected)
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			assert.Equal(t, "been closed", been.Closed().String())
		})
	})

	t.Run("been disconnected", func(t *testing.T) {
		t.Run("matches a connection that has been disconnected", func(t *testing.T) {
			fixture.AssertMatch(t, been.Disconnected(), disconnected)
		})

		t.Run("matches a connection that has been closed or reconnected, since it was disconnected first",
			func(t *testing.T) {
				fixture.AssertMatch(t, been.Disconnected(), closed)
				fixture.AssertMatch(t, been.Disconnected(), reconnected)
			})

		t.Run("does not match a connection that has not been disconnected", func(t *testing.T) {
			fixture.AssertNoMatch(t, been.Disconnected(), nothingHappened)
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			assert.Equal(t, "been disconnected", been.Disconnected().String())
		})
	})

	t.Run("been reconnected", func(t *testing.T) {
		t.Run("matches a connection that has been reconnected", func(t *testing.T) {
			fixture.AssertMatch(t, been.Reconnected(), reconnected)
		})

		t.Run("does not match a connection that has not been reconnected", func(t *testing.T) {
			fixture.AssertNoMatch(t, been.Reconnected(), nothingHappened)
			fixture.AssertNoMatch(t, been.Reconnected(), disconnected)
			fixture.AssertNoMatch(t, been.Reconnected(), closed)
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			assert.Equal(t, "been reconnected", been.Reconnected().String())
		})
	})

	t.Run("does not match a nil connection", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Closed(), nil)
		fixture.AssertNoMatch(t, been.Disconnected(), nil)
		fixture.AssertNoMatch(t, been.Reconnected(), nil)
	})

	t.Run("implements the resolution interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Resolution[*ability.Connection])(nil), been.Closed())
		assert.Implements(t, (*screenplay.Resolution[*ability.Connection])(nil), been.Disconnected())
		assert.Implements(t, (*screenplay.Resolution[*ability.Connection])(nil), been.Reconnected())
	})
}
