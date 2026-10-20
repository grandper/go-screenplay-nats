package been_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/resolution/testdata"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestClosedResolution(t *testing.T) {
	t.Run("should match if the connection is closed", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
		nc.Close()
		time.Sleep(100 * time.Millisecond)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS)
		connection, err := question.NATSConnection().AnsweredBy(adam)
		require.NoError(t, err)

		testdata.AssertMatch(t, been.Closed(), connection)
	})

	t.Run("fails when the connection is open", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS)
		connection, err := question.NATSConnection().AnsweredBy(adam)
		require.NoError(t, err)

		testdata.AssertNoMatch(t, been.Closed(), connection)
		testdata.AssertNoMatch(t, been.Closed(), nil)
	})

	t.Run("returns an error when the value is of the wrong type", func(t *testing.T) {
		testdata.AssertMatcherFails(t, been.Closed(), 2)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been closed", been.Closed().String())
	})
}
