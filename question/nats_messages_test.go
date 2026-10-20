package question_test

import (
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/question/number"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action/publish"
	"github.com/grandper/go-screenplay-nats/v2/action/receive"
	"github.com/grandper/go-screenplay-nats/v2/action/subscribe"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
)

const (
	timeout      = 2 * time.Second
	tickInterval = 100 * time.Millisecond
)

func TestNATSMessages(t *testing.T) {
	const (
		orders   = "orders"
		invoices = "invoices"
	)

	t.Run("answers with the messages received by the actor", func(t *testing.T) {
		nc := fixture.CreateEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		require.NoError(t, adam.AttemptsTo(
			subscribe.To(orders),
			subscribe.To(invoices),
			publish.TheMessage().To(orders).Containing([]byte("1")),
			receive.TheNextMessageFrom(orders),
			publish.TheMessage().To(invoices).Containing([]byte("2")),
			receive.TheNextMessageFrom(invoices),
			publish.TheMessage().To(orders).Containing([]byte("3")),
			receive.TheNextMessageFrom(orders),
		))

		t.Run("all of them, in the order they were received", func(t *testing.T) {
			messages, err := adam.AsksFor(question.NATSMessages())

			require.NoError(t, err)
			require.Len(t, messages, 3)
			assert.Equal(t, []byte("1"), messages[0].Data)
			assert.Equal(t, []byte("2"), messages[1].Data)
			assert.Equal(t, []byte("3"), messages[2].Data)
		})

		t.Run("those from a subject", func(t *testing.T) {
			messages, err := adam.AsksFor(question.NATSMessagesFrom(orders))

			require.NoError(t, err)
			require.Len(t, messages, 2)
			assert.Equal(t, []byte("1"), messages[0].Data)
			assert.Equal(t, []byte("3"), messages[1].Data)
		})

		t.Run("none from a subject nothing was received from", func(t *testing.T) {
			messages, err := adam.AsksFor(question.NATSMessagesFrom("payments"))

			require.NoError(t, err)
			assert.Empty(t, messages)
		})

		t.Run("can be composed with the core questions", func(t *testing.T) {
			require.NoError(t, adam.Should(see.The(number.Of(question.NATSMessages())).Should(be.EqualTo(3))))
			lastMessage, err := adam.AsksFor(last.Of(question.NATSMessagesFrom(orders)))
			require.NoError(t, err)
			assert.Equal(t, []byte("3"), lastMessage.Data)
		})
	})

	t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(question.NATSMessages())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = screenplay.ActorNamed("Adam").AsksFor(question.NATSMessagesFrom(orders))
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "NATS messages", question.NATSMessages().String())
		assert.Equal(t, `NATS messages from "orders"`, question.NATSMessagesFrom(orders).String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[[]*nats.Msg])(nil), question.NATSMessages())
		assert.Implements(t, (*screenplay.Question[[]*nats.Msg])(nil), question.NATSMessagesFrom(orders))
	})
}
