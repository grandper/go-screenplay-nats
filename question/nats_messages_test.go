package question_test

import (
	"testing"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay-nats/action"
	"github.com/grandper/go-screenplay-nats/fixture"
	"github.com/grandper/go-screenplay-nats/question"
	"github.com/grandper/go-screenplay-nats/resolution/been"
	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNATSMessages(t *testing.T) {
	nc := fixture.CreateEmbeddedNATSServer(t)
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
	const subject1 = "test-subject-1"
	const subject2 = "test-subject-2"
	content1 := []byte("Hello, NATS!")
	content2 := []byte("Hello, NATS again!")

	t.Run("returns all messages", func(t *testing.T) {
		require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject1)))
		require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject2)))

		require.NoError(t, adam.AttemptsTo(action.Publish(ability.Message().To(subject1).Containing(content1))))
		require.NoError(t, adam.AttemptsTo(action.ReceiveFrom(subject1)))
		require.NoError(t, adam.AttemptsTo(action.Publish(ability.Message().To(subject2).Containing(content2))))
		require.NoError(t, adam.AttemptsTo(action.ReceiveFrom(subject2)))

		t.Run("returns all messages when no subject is set", func(t *testing.T) {
			answer, err := adam.AsksFor(question.NATSMessages())
			require.NoError(t, err)
			msg := answer.([]*nats.Msg)
			require.Len(t, msg, 2)
		})

		t.Run("returns all messages for a specific subject", func(t *testing.T) {
			answer, err := adam.AsksFor(question.NATSMessages().From(subject1))
			require.NoError(t, err)
			msg := answer.([]*nats.Msg)
			require.Len(t, msg, 1)
		})
	})

	t.Run("fails to get the NATS messages if the actor does not have the ability CommunicateWithNATS",
		func(t *testing.T) {
			bob := screenplay.ActorNamed("Bob")
			assert.Error(t, bob.AttemptsTo(see.The(question.NATSMessages()).Have(been.Received())))
		},
	)

	t.Run("implements the stringer interface", func(t *testing.T) {
		body := question.NATSMessages()
		assert.Equal(t, "NATS messages", body.String())

		body = question.NATSMessages().From(subject1)
		assert.Equal(t, "NATS messages from test-subject-1", body.String())
	})
}
