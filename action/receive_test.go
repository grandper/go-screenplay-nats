package action_test

import (
	"testing"
	"time"

	act "github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/question"
	"github.com/grandper/go-screenplay-nats/v2/resolution/been"
)

func TestReceive(t *testing.T) {
	const subject = "test.subject"
	content := []byte("hello")
	const waitInterval = 100 * time.Millisecond

	nc := fixture.CreateEmbeddedNATSServer(t)
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))

	t.Run("can receive messages", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			require.NoError(t, adam.AttemptsTo(action.SubscribeTo(subject)))
			time.Sleep(waitInterval)
			require.NoError(t, nc.Publish(subject, content))
			require.NoError(t, adam.AttemptsTo(action.ReceiveFrom(subject)))
			require.NoError(
				t,
				adam.Should(act.Eventually(see.The(question.NATSMessages().From(subject)).Have(been.Received()))),
			)
		})

		t.Run("fail the connection is nil", func(t *testing.T) {
			bob := screenplay.ActorNamed("Bob").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, bob.AttemptsTo(action.ReceiveFrom(subject)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fail when no subscription exists for the subject", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.ReceiveFrom("unknown.subject")), ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the actor has not the ability to communicate with NATS", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.ReceiveFrom(subject)),
				screenplay.ErrActorMissingAbility)
		})

		t.Run("fails when the subject is not provided", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.ReceiveFrom("")), ability.ErrSubjectIsEmpty)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.ReceiveFrom(subject)
		assert.Equal(t, "receive a message from test.subject", action1.String())
	})
}
