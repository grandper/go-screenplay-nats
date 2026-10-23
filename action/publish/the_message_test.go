package publish_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/ability/communicate"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/publish"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
)

func TestTheMessage(t *testing.T) {
	const (
		subject = "orders"
		text    = "Hello, NATS!"
		timeout = 2 * time.Second
	)
	binary := []byte{0x00, 0x01, 0x02}

	t.Run("publishes a message", func(t *testing.T) {
		t.Run("with a text", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)

			require.NoError(t, adam.AttemptsTo(publish.TheMessage(text).To(subject)))

			msg, err := sub.NextMsg(timeout)
			require.NoError(t, err)
			assert.Equal(t, subject, msg.Subject)
			assert.Equal(t, []byte(text), msg.Data)
			assert.Empty(t, msg.Reply)
			assert.Empty(t, msg.Header)
		})

		t.Run("with binary content", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)

			require.NoError(t, adam.AttemptsTo(publish.TheBinaryMessage(binary).To(subject)))

			msg, err := sub.NextMsg(timeout)
			require.NoError(t, err)
			assert.Equal(t, binary, msg.Data)
		})

		t.Run("with no content", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)

			require.NoError(t, adam.AttemptsTo(publish.AnEmptyMessage().To(subject)))

			msg, err := sub.NextMsg(timeout)
			require.NoError(t, err)
			assert.Empty(t, msg.Data)
		})

		t.Run("with headers, expecting a reply", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)

			require.NoError(t, adam.AttemptsTo(
				publish.TheMessage(text).
					To(subject).
					WithTheHeader("Key", "value").
					WithTheHeaders("Other", "1", "Other", "2").
					ExpectingAReplyOn("inbox"),
			))

			msg, err := sub.NextMsg(timeout)
			require.NoError(t, err)
			assert.Equal(t, "value", msg.Header.Get("Key"))
			assert.Equal(t, []string{"1", "2"}, msg.Header.Values("Other"))
			assert.Equal(t, "inbox", msg.Reply)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			require.ErrorIs(t, adam.AttemptsTo(publish.TheMessage(text)), action.ErrSubjectIsEmpty)
			require.ErrorIs(t, adam.AttemptsTo(publish.TheMessage(text).To("")), action.ErrSubjectIsEmpty)
		})

		t.Run("fails when the headers are written wrongly", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			err := adam.AttemptsTo(publish.TheMessage(text).To(subject).WithTheHeader("", "value"))
			require.ErrorIs(t, err, action.ErrInvalidHeaders)
			require.EqualError(t, err, "invalid headers: a key cannot be empty")

			err = adam.AttemptsTo(publish.TheMessage(text).To(subject).WithTheHeaders("Key"))
			require.ErrorIs(t, err, action.ErrInvalidHeaders)
			require.EqualError(t, err, "invalid headers: a key is missing its value")

			err = adam.AttemptsTo(publish.TheMessage(text).To(subject).WithTheHeaders())
			require.ErrorIs(t, err, action.ErrInvalidHeaders)
			require.EqualError(t, err, "invalid headers: at least one key and its value must be given")
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			require.ErrorIs(t, adam.AttemptsTo(publish.TheMessage(text).To(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(publish.TheMessage(text).To(subject))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := internalfixture.NarratedActor("Adam", recorder, fixture.ConnectToAnEmbeddedNATSServer(t))

		require.NoError(t, adam.AttemptsTo(publish.TheMessage(text).To(subject)))

		assert.Equal(t, `Adam attempts to publish the message "Hello, NATS!" to "orders"`, recorder.Messages()[0])
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `publish the message "Hello, NATS!"`, publish.TheMessage(text).String())
		assert.Equal(t, `publish the message "Hello, NATS!" to "orders"`, publish.TheMessage(text).To(subject).String())
		assert.Equal(
			t,
			`publish the binary message "\x00\x01\x02" to "orders"`,
			publish.TheBinaryMessage(binary).To(subject).String(),
		)
		assert.Equal(t, `publish an empty message to "orders"`, publish.AnEmptyMessage().To(subject).String())
		assert.Equal(
			t,
			`publish the message "Hello, NATS!" to "orders" with the header "Key" = "value"`+
				` with the headers "B" = "2", "A" = "1", "B" = "3" expecting a reply on "inbox"`,
			publish.TheMessage(text).
				To(subject).
				WithTheHeader("Key", "value").
				WithTheHeaders("B", "2", "A", "1", "B", "3").
				ExpectingAReplyOn("inbox").
				String(),
		)
		assert.Equal(
			t,
			`publish the message "Hello, NATS!" with the header "Key" = "value" to "orders"`,
			publish.TheMessage(text).WithTheHeader("Key", "value").To(subject).String(),
		)
	})
}
