package messages_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
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
	"github.com/grandper/go-screenplay-nats/v2/question/messages"
)

func TestFromAnySubject(t *testing.T) {
	const (
		orders   = "orders"
		invoices = "invoices"
	)

	t.Run("answers with every message the actor received", func(t *testing.T) {
		nc := fixture.ConnectToAnEmbeddedNATSServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
		require.NoError(t, adam.AttemptsTo(
			subscribe.To(orders),
			subscribe.To(invoices),
			publish.TheMessage("1").To(orders),
			receive.TheNextMessageFrom(orders),
			publish.TheMessage("2").To(invoices),
			receive.TheNextMessageFrom(invoices),
			publish.TheMessage("3").To(orders),
			receive.TheNextMessageFrom(orders),
		))

		t.Run("in the order they were received", func(t *testing.T) {
			received, err := adam.AsksFor(messages.FromAnySubject())

			require.NoError(t, err)
			require.Len(t, received, 3)
			assert.Equal(t, []byte("1"), received[0].Data)
			assert.Equal(t, []byte("2"), received[1].Data)
			assert.Equal(t, []byte("3"), received[2].Data)
		})

		t.Run("can be composed with the core questions", func(t *testing.T) {
			require.NoError(t, adam.Should(see.The(number.Of(messages.FromAnySubject())).Should(be.EqualTo(3))))
		})
	})

	t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(messages.FromAnySubject())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "NATS messages from any subject", messages.FromAnySubject().String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[[]*nats.Msg])(nil), messages.FromAnySubject())
	})
}
