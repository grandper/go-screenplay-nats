package been_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestReceivedResolution(t *testing.T) {
	t.Run("should match if messages have been received", func(t *testing.T) {
		fixture.AssertMatch(t, been.Received(), []*nats.Msg{
			{Subject: "foo"},
		})
	})

	t.Run("should match if messages have been received from a specific subject", func(t *testing.T) {
		fixture.AssertMatch(t, been.Received().From("bar"), []*nats.Msg{
			{Subject: "foo"},
			{Subject: "bar"},
		})
	})

	t.Run("fails when messages have not been received from a specific subject", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Received().From("bar"), []*nats.Msg{
			{Subject: "foo"},
		})
	})

	t.Run("fails when no messages have been received", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Received(), []*nats.Msg{})
	})

	t.Run("fails when no messages have been received (nil)", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Received(), nil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been received", been.Received().String())
	})
}
