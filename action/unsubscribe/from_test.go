package unsubscribe_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/subscribe"
	"github.com/grandper/go-screenplay-nats/v2/action/unsubscribe"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
)

func TestFrom(t *testing.T) {
	const subject = "orders"

	t.Run("unsubscribes from a subject", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.NoError(t, adam.AttemptsTo(subscribe.To(subject)))
			require.Equal(t, 1, nc.NumSubscriptions())

			require.NoError(t, adam.AttemptsTo(unsubscribe.From(subject)))

			assert.Equal(t, 0, nc.NumSubscriptions())
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(unsubscribe.From("")), action.ErrSubjectIsEmpty)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(unsubscribe.From(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when there is no subscription to the subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.ErrorIs(t, adam.AttemptsTo(unsubscribe.From(subject)), ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(unsubscribe.From(subject))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := internalfixture.NarratedActor("Adam", recorder, fixture.CreateEmbeddedNATSServer(t))
		require.NoError(t, adam.AttemptsTo(subscribe.To(subject)))

		require.NoError(t, adam.AttemptsTo(unsubscribe.From(subject)))

		assert.Equal(t, `Adam attempts to unsubscribe from "orders"`, recorder.Messages()[2])
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `unsubscribe from "orders"`, unsubscribe.From(subject).String())
	})
}
