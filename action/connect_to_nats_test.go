package action_test

import (
	"testing"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay-nats/action"
	"github.com/grandper/go-screenplay-nats/fixture"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectToNATS(t *testing.T) {
	natsServer := fixture.RunNATSServer(t)

	url := natsServer.ClientURL()
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS())

	t.Run("connects to NATS server", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			require.NoError(t, adam.AttemptsTo(action.ConnectToNATS().At(url)))
			// Verify the ability was updated in the actor by performing a subsequent action.
			require.NoError(t, adam.AttemptsTo(action.SubscribeTo("test.subject")))
		})

		t.Run("fail when the URL is empty", func(t *testing.T) {
			require.ErrorIs(t, adam.AttemptsTo(action.ConnectToNATS()), action.ErrURLIsEmpty)
			require.ErrorIs(t, adam.AttemptsTo(action.ConnectToNATS().At("")), action.ErrURLIsEmpty)
		})

		t.Run("fails when the connection cannot be established", func(t *testing.T) {
			const invalidURL = "nats://invalid:4222"
			eve := screenplay.ActorNamed("Eve").WhoCan(ability.CommunicateWithNATS())
			require.ErrorIs(t, eve.AttemptsTo(action.ConnectToNATS().At(invalidURL)),
				action.ErrFailedToConnect)
		})

		t.Run("fails when the actor has not the ability to communicate with NATS", func(t *testing.T) {
			require.ErrorIs(t, screenplay.ActorNamed("Bob").AttemptsTo(action.ConnectToNATS().At(url)),
				screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		action1 := action.ConnectToNATS().At("nats://localhost:4222")
		assert.Equal(t, "Connect to NATS server at nats://localhost:4222", action1.String())
	})
}
