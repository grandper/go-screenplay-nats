package action_test

import (
	"sync"
	"testing"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay-nats/action"
	"github.com/grandper/go-screenplay-nats/fixture"
	"github.com/grandper/go-screenplay-nats/question"
	"github.com/grandper/go-screenplay-nats/resolution/been"
	act "github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	const subject = "test.subject"
	var content = []byte("Hello, NATS!")

	aMessage := ability.Message().To(subject).Containing(content)

	nc := fixture.CreateEmbeddedNATSServer(t)
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))

	t.Run("can send a request", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				assert.NoError(t, msg.Respond([]byte("ok")))
			})
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()
			require.NoError(t, adam.AttemptsTo(action.Request(aMessage)))
		})

		t.Run("stores the response for later retrieval", func(t *testing.T) {
			alice := screenplay.ActorNamed("Alice").WhoCan(ability.CommunicateWithNATS().Using(nc))

			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				assert.NoError(t, msg.Respond([]byte("response")))
			})
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()

			require.NoError(t, alice.AttemptsTo(action.Request(aMessage)))
			require.NoError(t, alice.Should(act.Eventually(see.The(
				question.NATSMessages().From(subject),
			).Have(been.Received()))))
		})

		t.Run("fails when the response is not fast enough", func(t *testing.T) {
			var mutex sync.Mutex
			messageReceived := false
			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				mutex.Lock()
				defer mutex.Unlock()
				messageReceived = true
			})
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()

			err = adam.AttemptsTo(action.Request(aMessage).For(20).Milliseconds())
			require.ErrorIs(t, err, nats.ErrTimeout)

			mutex.Lock()
			defer mutex.Unlock()
			assert.True(t, messageReceived)
		})

		t.Run("fail the connection is nil", func(t *testing.T) {
			bob := screenplay.ActorNamed("Bob").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, bob.AttemptsTo(action.Request(aMessage)), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fail if the message is nil", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.Request(nil)), ability.ErrMessageIsNil)
		})

		t.Run("fails when the actor has not the ability to communicate with NATS", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.Request(aMessage)),
				screenplay.ErrActorMissingAbility)
		})

		t.Run("fails when the subject is not provided", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.Request(ability.Message())), ability.ErrSubjectIsEmpty)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.Request(aMessage)
		assert.Equal(t, "send a request to test.subject for 5 seconds", action1.String())

		actionNil := action.Request(nil)
		assert.Equal(t, "send a nil request", actionNil.String())

		action2 := action.Request(aMessage).For(0).Second()
		assert.Equal(t, "send a request to test.subject for 0 second", action2.String())

		action3 := action.Request(aMessage).For(1).Second()
		assert.Equal(t, "send a request to test.subject for 1 second", action3.String())

		action4 := action.Request(aMessage).For(6).Seconds()
		assert.Equal(t, "send a request to test.subject for 6 seconds", action4.String())

		action5 := action.Request(aMessage).For(0).Millisecond()
		assert.Equal(t, "send a request to test.subject for 0 millisecond", action5.String())

		action6 := action.Request(aMessage).For(1).Millisecond()
		assert.Equal(t, "send a request to test.subject for 1 millisecond", action6.String())

		action7 := action.Request(aMessage).For(6).Milliseconds()
		assert.Equal(t, "send a request to test.subject for 6 milliseconds", action7.String())
	})
}
