package action_test

import (
	"testing"

	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
)

func TestFlush(t *testing.T) {
	const subject = "test.subject"
	var content = []byte("Hello, NATS!")

	aMessage := ability.Message().To(subject).Containing(content)

	nc := fixture.CreateEmbeddedNATSServer(t)
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))

	t.Run("can flush the connection", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			require.NoError(t, adam.AttemptsTo(action.Publish(aMessage)))
			require.NoError(t, adam.AttemptsTo(action.Flush()))
		})

		t.Run("fail the connection is nil", func(t *testing.T) {
			bob := screenplay.ActorNamed("Bob").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, bob.AttemptsTo(action.Flush()), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the actor has not the ability to communicate with NATS", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.Flush()),
				screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Flush()
		assert.Equal(t, "flush the NATS connection", action1.String())
	})
}
