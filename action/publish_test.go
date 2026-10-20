package action_test

import (
	"sync"
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

func TestPublish(t *testing.T) {
	const subject = "test.subject"
	var content = []byte("Hello, NATS!")
	const timeout = 2 * time.Second
	const waitInterval = 100 * time.Millisecond

	aMessage := ability.Message().To(subject).Containing(content)

	nc := fixture.CreateEmbeddedNATSServer(t)
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))

	t.Run("can publish messages", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			var mutex sync.Mutex
			msgReceived := false
			sub, err := nc.Subscribe(subject, func(_ *nats.Msg) {
				mutex.Lock()
				defer mutex.Unlock()
				msgReceived = true
			})
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()
			require.NoError(t, adam.AttemptsTo(action.Publish(aMessage)))
			assert.Eventually(t, func() bool {
				mutex.Lock()
				defer mutex.Unlock()
				return msgReceived
			}, timeout, waitInterval)
		})

		t.Run("fail the connection is nil", func(t *testing.T) {
			bob := screenplay.ActorNamed("Bob").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, bob.AttemptsTo(action.Publish(aMessage)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fail if the message is nil", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.Publish(nil)), ability.ErrMessageIsNil)
		})

		t.Run("fails when the actor has not the ability to communicate with NATS", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.Publish(aMessage)),
				screenplay.ErrActorMissingAbility)
		})

		t.Run("fails when the subject is not provided", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.Publish(ability.Message())), action.ErrSubjectIsEmpty)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Publish(aMessage)
		assert.Equal(t, "publish a message to test.subject", action1.String())

		action2 := action.Publish(nil)
		assert.Equal(t, "publish a nil message", action2.String())
	})
}
