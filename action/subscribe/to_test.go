package subscribe_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/ability/communicate"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/receive"
	"github.com/grandper/go-screenplay-nats/v2/action/subscribe"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
)

func TestTo(t *testing.T) {
	const (
		subject = "orders"
		queue   = "workers"
	)

	t.Run("subscribes to a subject", func(t *testing.T) {
		t.Run("synchronously by default", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			require.NoError(t, adam.AttemptsTo(subscribe.To(subject)))

			assert.Equal(t, 1, nc.NumSubscriptions())
			require.NoError(t, nc.Publish(subject, []byte("hello")))
			require.NoError(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)))
		})

		t.Run("synchronously", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			require.NoError(t, adam.AttemptsTo(subscribe.To(subject).Synchronously()))

			require.NoError(t, nc.Publish(subject, []byte("hello")))
			require.NoError(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)))
		})

		t.Run("asynchronously", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			require.NoError(t, adam.AttemptsTo(subscribe.To(subject).Asynchronously()))

			// The next message cannot be received by hand from an asynchronous subscription.
			require.NoError(t, nc.Publish(subject, []byte("hello")))
			require.ErrorIs(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)), nats.ErrSyncSubRequired)
		})

		t.Run("using a queue", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			require.NoError(t, adam.AttemptsTo(subscribe.To(subject).InTheQueueGroup(queue)))

			assert.Equal(t, 1, nc.NumSubscriptions())
		})

		t.Run("using a queue asynchronously", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			require.NoError(t, adam.AttemptsTo(subscribe.To(subject).InTheQueueGroup(queue).Asynchronously()))

			require.NoError(t, nc.Publish(subject, []byte("hello")))
			require.ErrorIs(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)), nats.ErrSyncSubRequired)
		})

		t.Run("and unsubscribes after a number of messages", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			require.NoError(t, adam.AttemptsTo(subscribe.To(subject).AndUnsubscribeAfter(1).Message()))

			require.NoError(t, nc.Publish(subject, []byte("hello")))
			require.NoError(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)))
			require.ErrorIs(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject)), nats.ErrMaxMessages)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			require.ErrorIs(t, adam.AttemptsTo(subscribe.To("")), action.ErrSubjectIsEmpty)
		})

		t.Run("fails when the queue is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			require.ErrorIs(t, adam.AttemptsTo(subscribe.To(subject).InTheQueueGroup("")), action.ErrQueueIsEmpty)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			require.ErrorIs(t, adam.AttemptsTo(subscribe.To(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the subscription already exists", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			require.NoError(t, adam.AttemptsTo(subscribe.To(subject)))
			require.ErrorIs(t, adam.AttemptsTo(subscribe.To(subject)), ability.ErrSubscriptionAlreadyExists)
		})

		t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(subscribe.To(subject))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("subscribes to the subjects starting with a prefix", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

			require.NoError(t, adam.AttemptsTo(subscribe.ToTheSubjectsStartingWith(subject)))

			require.NoError(t, nc.Publish(subject+".created", []byte("hello")))
			require.NoError(t, adam.AttemptsTo(receive.TheNextMessageFrom(subject+".>")))
		})

		t.Run("fails when the prefix is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			require.ErrorIs(t, adam.AttemptsTo(subscribe.ToTheSubjectsStartingWith("")), action.ErrSubjectIsEmpty)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := internalfixture.NarratedActor("Adam", recorder, fixture.ConnectToAnEmbeddedNATSServer(t))

		require.NoError(
			t,
			adam.AttemptsTo(subscribe.To(subject).InTheQueueGroup(queue).AndUnsubscribeAfter(2).Messages()),
		)

		assert.Equal(
			t,
			`Adam attempts to subscribe to "orders" in the queue group "workers" and unsubscribe after 2 messages`,
			recorder.Messages()[0],
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `subscribe to "orders"`, subscribe.To(subject).String())
		assert.Equal(t, `subscribe to "orders" synchronously`, subscribe.To(subject).Synchronously().String())
		assert.Equal(t, `subscribe to "orders" asynchronously`, subscribe.To(subject).Asynchronously().String())
		assert.Equal(
			t,
			`subscribe to "orders" in the queue group "workers"`,
			subscribe.To(subject).InTheQueueGroup(queue).String(),
		)
		assert.Equal(
			t,
			`subscribe to "orders" asynchronously in the queue group "workers"`,
			subscribe.To(subject).Asynchronously().InTheQueueGroup(queue).String(),
		)
		assert.Equal(
			t,
			`subscribe to "orders" and unsubscribe after 4 messages`,
			subscribe.To(subject).AndUnsubscribeAfter(4).Messages().String(),
		)
		assert.Equal(
			t,
			`subscribe to "orders" and unsubscribe after 1 message`,
			subscribe.To(subject).AndUnsubscribeAfter(1).Message().String(),
		)
		assert.Equal(
			t,
			`subscribe to the subjects starting with "orders"`,
			subscribe.ToTheSubjectsStartingWith(subject).String(),
		)
	})
}
