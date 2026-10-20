package action_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
)

func TestSubscribeTo(t *testing.T) {
	const (
		subject      = "test.subject"
		queue        = "test.queue"
		timeout      = 2 * time.Second
		waitInterval = 100 * time.Millisecond
	)

	t.Run("can subscribe to subject synchronously", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		t.Run("successfully", func(t *testing.T) {
			require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject)))
			assert.Eventually(t, func() bool {
				return nc.NumSubscriptions() >= 1
			}, timeout, waitInterval)
		})

		t.Run("fail the connection is nil", func(t *testing.T) {
			bob := screenplay.ActorNamed("Bob").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, bob.AttemptsTo(action.SubscribeTo(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fail when the subscription already exists for the subject", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.SubscribeTo(subject)), ability.ErrSubscriptionAlreadyExists)
		})

		t.Run("fails when the actor has not the ability to communicate with NATS", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.SubscribeTo(subject)),
				screenplay.ErrActorMissingAbility)
		})

		t.Run("fails when the subject is not provided", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.SubscribeTo("")), action.ErrSubjectIsEmpty)
		})
	})

	t.Run("can subscribe to subject starting with a given prefix", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		require.NoError(t, adam.AttemptsTo(action.SubscribeToTheSubjectStartingWith(subject)))
		assert.Eventually(t, func() bool {
			return nc.NumSubscriptions() >= 1
		}, timeout, waitInterval)
	})

	t.Run("can subscribe to subject using a queue", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject).UsingTheQueue(queue)))
		assert.Eventually(t, func() bool {
			return nc.NumSubscriptions() >= 1
		}, timeout, waitInterval)
	})

	t.Run("can subscribe to subject using a queue asynchronously", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject).UsingTheQueue(queue).Asynchronously()))
		assert.Eventually(t, func() bool {
			return nc.NumSubscriptions() >= 1
		}, timeout, waitInterval)
		// Verify the subscription is truly async: NextMsg is illegal on async subscriptions.
		require.NoError(t, nc.Publish(subject, []byte("hello")))
		require.ErrorIs(t, adam.AttemptsTo(action.ReceiveFrom(subject)), nats.ErrSyncSubRequired)
	})

	t.Run("can subscribe to subject asynchronously", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject).Asynchronously()))
		assert.Eventually(t, func() bool {
			return nc.NumSubscriptions() >= 1
		}, timeout, waitInterval)
		// Verify the subscription is truly async: NextMsg is illegal on async subscriptions.
		require.NoError(t, nc.Publish(subject, []byte("hello")))
		require.ErrorIs(t, adam.AttemptsTo(action.ReceiveFrom(subject)), nats.ErrSyncSubRequired)
	})

	t.Run("can subscribe for a given number of messages", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject).AndUnsubscribeAfter(5).Messages()))
		assert.Eventually(t, func() bool {
			return nc.NumSubscriptions() >= 1
		}, timeout, waitInterval)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.SubscribeTo(subject)
		assert.Equal(t, "subscribe to NATS subject test.subject synchronously", action1.String())

		action2 := action.SubscribeTo(subject).Synchronously()
		assert.Equal(t, "subscribe to NATS subject test.subject synchronously", action2.String())

		action3 := action.SubscribeTo(subject).Asynchronously()
		assert.Equal(t, "subscribe to NATS subject test.subject asynchronously", action3.String())

		action4 := action.SubscribeTo(subject).UsingTheQueue(queue)
		assert.Equal(t, "subscribe to NATS subject test.subject synchronously using queue test.queue", action4.String())

		action5 := action.SubscribeTo(subject).AndUnsubscribeAfter(4).Messages()
		assert.Equal(
			t,
			"subscribe to NATS subject test.subject synchronously and unsubscribe after 4 messages",
			action5.String(),
		)

		action6 := action.SubscribeTo(subject).AndUnsubscribeAfter(1).Message()
		assert.Equal(
			t,
			"subscribe to NATS subject test.subject synchronously and unsubscribe after 1 message",
			action6.String(),
		)
	})
}
