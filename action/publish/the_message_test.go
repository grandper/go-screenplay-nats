package publish_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/publish"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
)

func TestTheMessage(t *testing.T) {
	const (
		subject = "orders"
		timeout = 2 * time.Second
	)
	content := []byte("Hello, NATS!")

	t.Run("publishes a message", func(t *testing.T) {
		t.Run("to a subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)

			require.NoError(t, adam.AttemptsTo(publish.TheMessage().To(subject)))

			msg, err := sub.NextMsg(timeout)
			require.NoError(t, err)
			assert.Equal(t, subject, msg.Subject)
			assert.Empty(t, msg.Data)
			assert.Empty(t, msg.Reply)
		})

		t.Run("containing a payload, with headers, to be replied to", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)

			require.NoError(t, adam.AttemptsTo(
				publish.TheMessage().
					To(subject).
					Containing(content).
					WithTheHeader("Key", "value").
					WithTheHeaders(nats.Header{"Other": []string{"1", "2"}}).
					ToBeRepliedTo("inbox"),
			))

			msg, err := sub.NextMsg(timeout)
			require.NoError(t, err)
			assert.Equal(t, content, msg.Data)
			assert.Equal(t, "value", msg.Header.Get("Key"))
			assert.Equal(t, []string{"1", "2"}, msg.Header.Values("Other"))
			assert.Equal(t, "inbox", msg.Reply)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(publish.TheMessage()), action.ErrSubjectIsEmpty)
			require.ErrorIs(t, adam.AttemptsTo(publish.TheMessage().To("")), action.ErrSubjectIsEmpty)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(publish.TheMessage().To(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(publish.TheMessage().To(subject))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := internalfixture.NarratedActor("Adam", recorder, fixture.CreateEmbeddedNATSServer(t))

		require.NoError(t, adam.AttemptsTo(publish.TheMessage().To(subject).Containing(content)))

		assert.Equal(
			t,
			`Adam attempts to publish the message to "orders" containing "Hello, NATS!"`,
			recorder.Messages()[0],
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "publish the message", publish.TheMessage().String())
		assert.Equal(t, `publish the message to "orders"`, publish.TheMessage().To(subject).String())
		assert.Equal(
			t,
			`publish the message to "orders" containing "Hello, NATS!" with the header "Key" = "value"`+
				` with the headers "A" = "1", "B" = "2", "B" = "3" to be replied to "inbox"`,
			publish.TheMessage().
				To(subject).
				Containing(content).
				WithTheHeader("Key", "value").
				WithTheHeaders(nats.Header{"B": []string{"2", "3"}, "A": []string{"1"}}).
				ToBeRepliedTo("inbox").
				String(),
		)
		assert.Equal(
			t,
			`publish the message containing "Hello, NATS!" to "orders"`,
			publish.TheMessage().Containing(content).To(subject).String(),
		)
	})
}
