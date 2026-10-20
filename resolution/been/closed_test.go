package been_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/stretchr/testify/assert"

	questionfixture "github.com/grandper/go-screenplay-nats/v2/question/fixture"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestClosed(t *testing.T) {
	t.Run("matches a connection that has been closed", func(t *testing.T) {
		fixture.AssertMatch(t, been.Closed(), questionfixture.CreateFakeClosedConnection())
	})

	t.Run("does not match a connection that has not been closed", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Closed(), questionfixture.CreateFakeConnection())
	})

	t.Run("does not match a nil connection", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Closed(), nil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been closed", been.Closed().String())
	})
}
