package send_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/send"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question/messages"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestTheRequest(t *testing.T) {
	const (
		subject = "orders"
		text    = "ping"
	)
	binary := []byte{0x00, 0x01}

	t.Run("sends a request", func(t *testing.T) {
		t.Run("with a text, and stores the reply", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				assert.Equal(t, []byte(text), msg.Data)
				assert.Equal(t, "value", msg.Header.Get("Key"))
				assert.Equal(t, "1", msg.Header.Get("Other"))
				assert.NoError(t, msg.Respond([]byte("pong")))
			})
			require.NoError(t, err)
			defer func() { assert.NoError(t, sub.Unsubscribe()) }()

			require.NoError(t, adam.AttemptsTo(
				send.TheRequest(text).To(subject).WithTheHeader("Key", "value").WithTheHeaders("Other", "1"),
			))

			require.NoError(t, adam.Should(see.The(messages.FromTheSubject(subject)).ShouldHave(been.Received())))
			replies, err := adam.AsksFor(messages.FromTheSubject(subject))
			require.NoError(t, err)
			require.Len(t, replies, 1)
			assert.Equal(t, []byte("pong"), replies[0].Data)
		})

		t.Run("with binary content", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				assert.Equal(t, binary, msg.Data)
				assert.NoError(t, msg.Respond(nil))
			})
			require.NoError(t, err)
			defer func() { assert.NoError(t, sub.Unsubscribe()) }()

			require.NoError(t, adam.AttemptsTo(send.TheBinaryRequest(binary).To(subject)))
		})

		t.Run("with no content", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				assert.Empty(t, msg.Data)
				assert.NoError(t, msg.Respond(nil))
			})
			require.NoError(t, err)
			defer func() { assert.NoError(t, sub.Unsubscribe()) }()

			require.NoError(t, adam.AttemptsTo(send.AnEmptyRequest().To(subject)))
		})

		t.Run("fails when the reply does not come in time", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
			sub, err := nc.Subscribe(subject, func(_ *nats.Msg) {})
			require.NoError(t, err)
			defer func() { assert.NoError(t, sub.Unsubscribe()) }()

			err = adam.AttemptsTo(send.TheRequest(text).To(subject).TimingOutAfter(20).Milliseconds())

			require.ErrorIs(t, err, nats.ErrTimeout)
		})

		t.Run("fails when nobody answers", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
			require.ErrorIs(t, adam.AttemptsTo(send.TheRequest(text).To(subject)), nats.ErrNoResponders)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(send.TheRequest(text)), action.ErrSubjectIsEmpty)
		})

		t.Run("fails when the headers are written wrongly", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			err := adam.AttemptsTo(send.TheRequest(text).To(subject).WithTheHeaders("Key"))
			require.ErrorIs(t, err, action.ErrInvalidHeaders)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, adam.AttemptsTo(send.TheRequest(text).To(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(send.TheRequest(text).To(subject))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		nc := fixture.ConnectToAnEmbeddedNATSServer(t)
		adam := internalfixture.NarratedActor("Adam", recorder, nc)
		sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
			assert.NoError(t, msg.Respond([]byte("pong")))
		})
		require.NoError(t, err)
		defer func() { assert.NoError(t, sub.Unsubscribe()) }()

		require.NoError(t, adam.AttemptsTo(send.TheRequest(text).To(subject).TimingOutAfter(1).Second()))

		assert.Equal(
			t,
			`Adam attempts to send the request "ping" to "orders" timing out after 1 second`,
			recorder.Messages()[0],
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `send the request "ping"`, send.TheRequest(text).String())
		assert.Equal(t, `send the request "ping" to "orders"`, send.TheRequest(text).To(subject).String())
		assert.Equal(
			t,
			`send the binary request "\x00\x01" to "orders"`,
			send.TheBinaryRequest(binary).To(subject).String(),
		)
		assert.Equal(t, `send an empty request to "orders"`, send.AnEmptyRequest().To(subject).String())
		assert.Equal(
			t,
			`send the request "ping" to "orders" with the header "Key" = "value"`+
				` with the headers "A" = "1" timing out after 500 milliseconds`,
			send.TheRequest(text).
				To(subject).
				WithTheHeader("Key", "value").
				WithTheHeaders("A", "1").
				TimingOutAfter(500).Milliseconds().
				String(),
		)
	})
}
