package were_test

import (
	"testing"

	"github.com/grandper/go-screenplay-nats/resolution/were"
	"github.com/grandper/go-screenplay/resolution/testdata"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestCReceivedResolution(t *testing.T) {
	t.Run("should match if messages were received", func(t *testing.T) {
		testdata.AssertMatch(t, were.Received(), []*nats.Msg{
			{Subject: "foo"},
		})
	})

	t.Run("should match if messages were received from a specific subject", func(t *testing.T) {
		testdata.AssertMatch(t, were.Received().From("bar"), []*nats.Msg{
			{Subject: "foo"},
			{Subject: "bar"},
		})
	})

	t.Run("fails when messages were not received from a specific subject", func(t *testing.T) {
		testdata.AssertNoMatch(t, were.Received().From("bar"), []*nats.Msg{
			{Subject: "foo"},
		})
	})

	t.Run("fails when no messages were received", func(t *testing.T) {
		testdata.AssertNoMatch(t, were.Received(), []*nats.Msg{})
	})

	t.Run("fails when no messages were received (nil)", func(t *testing.T) {
		testdata.AssertNoMatch(t, were.Received(), nil)
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, were.Received(), 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "were received", were.Received().String())
	})
}
