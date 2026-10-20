package receive_test

import (
	"context"
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/receive"
	"github.com/grandper/go-screenplay-nats/v2/action/subscribe"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
)

func TestTheNextMessageFrom(t *testing.T) {
	const subject = "orders"
	content := []byte("hello")

	t.Run("receives the next message", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.NoError(t, adam.AttemptsTo(subscribe.To(subject)))
			require.NoError(t, nc.Publish(subject, content))

			require.NoError(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)))

			messages, err := adam.AsksFor(question.NATSMessagesFrom(subject))
			require.NoError(t, err)
			require.Len(t, messages, 1)
			assert.Equal(t, content, messages[0].Data)
		})

		t.Run("waits at most the timeout of the actor", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			production := screenplay.NewProduction(screenplay.WithTimeout(50 * time.Millisecond))
			adam := production.SetTheStage(screenplay.CastOfStandardActors()).
				ActorNamed("Adam").
				WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.NoError(t, adam.AttemptsTo(subscribe.To(subject)))

			start := time.Now()
			err := adam.AttemptsTo(receive.TheNextMessageFrom(subject))

			require.ErrorIs(t, err, context.DeadlineExceeded)
			assert.Less(t, time.Since(start), time.Second)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(receive.TheNextMessageFrom("")), action.ErrSubjectIsEmpty)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when there is no subscription to the subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			require.ErrorIs(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)), ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(receive.TheNextMessageFrom(subject))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := internalfixture.NarratedActor("Adam", recorder, nc)
		require.NoError(t, adam.AttemptsTo(subscribe.To(subject)))
		require.NoError(t, nc.Publish(subject, content))

		require.NoError(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)))

		assert.Equal(t, `Adam attempts to receive the next message from "orders"`, recorder.Messages()[2])
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `receive the next message from "orders"`, receive.TheNextMessageFrom(subject).String())
	})
}
