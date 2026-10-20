package action_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
)

func TestUnsubscribeFrom(t *testing.T) {
	const subject = "test.subject"

	t.Run("can unsubscribe from a topic", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject)))
			require.Equal(t, 1, nc.NumSubscriptions())
			require.NoError(t, adam.AttemptsTo(action.UnsubscribeFrom(subject)))
			require.Equal(t, 0, nc.NumSubscriptions())
		})

		t.Run("fail the connection is nil", func(t *testing.T) {
			bob := screenplay.ActorNamed("Bob").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, bob.AttemptsTo(action.UnsubscribeFrom(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the actor has not the ability to communicate with NATS", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.UnsubscribeFrom(subject)),
				screenplay.ErrActorMissingAbility)
		})

		t.Run("fails when the subject is not provided", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.ErrorIs(t, adam.AttemptsTo(action.UnsubscribeFrom("")), ability.ErrSubjectIsEmpty)
		})

		t.Run("fails when no subscription exists for the subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.ErrorIs(
				t,
				adam.AttemptsTo(action.UnsubscribeFrom("non.existent.subject")),
				ability.ErrNoSubscriptionFound,
			)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.UnsubscribeFrom(subject)
		assert.Equal(t, "unsubscribe from NATS subject test.subject", action1.String())
	})
}
