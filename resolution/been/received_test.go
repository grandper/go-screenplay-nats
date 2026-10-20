package been_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/action/publish"
	"github.com/grandper/go-screenplay-nats/v2/action/receive"
	"github.com/grandper/go-screenplay-nats/v2/action/subscribe"
	natsfixture "github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestReceived(t *testing.T) {
	orders := &nats.Msg{Subject: "orders"}
	invoices := &nats.Msg{Subject: "invoices"}

	t.Run("matches when at least one message has been received", func(t *testing.T) {
		fixture.AssertMatch(t, been.Received(), []*nats.Msg{orders})
		fixture.AssertMatch(t, been.Received(), []*nats.Msg{orders, invoices})
	})

	t.Run("does not match when no message has been received", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Received(), []*nats.Msg{})
		fixture.AssertNoMatch(t, been.Received(), nil)
	})

	t.Run("matches when a message has been received from the subject", func(t *testing.T) {
		fixture.AssertMatch(t, been.ReceivedFrom("invoices"), []*nats.Msg{orders, invoices})
	})

	t.Run("does not match when no message has been received from the subject", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.ReceivedFrom("invoices"), []*nats.Msg{orders})
		fixture.AssertNoMatch(t, been.ReceivedFrom("invoices"), []*nats.Msg{nil})
		fixture.AssertNoMatch(t, been.ReceivedFrom("invoices"), nil)
	})

	t.Run("reads as a sentence with the question", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		nc := natsfixture.ConnectToAnEmbeddedNATSServer(t)
		adam := internalfixture.NarratedActor("Adam", recorder, nc)
		require.NoError(t, adam.AttemptsTo(subscribe.To("orders")))

		err := adam.Should(see.The(question.NATSMessagesFrom("orders")).ShouldHave(been.Received()))
		require.EqualError(
			t,
			err,
			`Adam failed to see that the NATS messages from "orders" should have been received: no message has been received`,
		)

		require.NoError(t, adam.AttemptsTo(
			publish.TheMessage().To("orders").Containing([]byte("Hello")),
			receive.TheNextMessageFrom("orders"),
		))
		require.NoError(t, adam.Should(see.The(question.NATSMessages()).ShouldHave(been.ReceivedFrom("orders"))))
		assert.Contains(
			t,
			recorder.Messages(),
			`Adam should see that the NATS messages should have been received from "orders"`,
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been received", been.Received().String())
		assert.Equal(t, `been received from "orders"`, been.ReceivedFrom("orders").String())
	})
}
