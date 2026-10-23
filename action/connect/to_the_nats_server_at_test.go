package connect_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/ability"
	"github.com/grandper/go-screenplay-nats/v2/ability/communicate"
	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/connect"
	"github.com/grandper/go-screenplay-nats/v2/action/flush"
	"github.com/grandper/go-screenplay-nats/v2/fixture"
	internalfixture "github.com/grandper/go-screenplay-nats/v2/internal/fixture"
)

func TestToTheNATSServerAt(t *testing.T) {
	t.Run("connects to the NATS server", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			natsServer := fixture.RunAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())

			require.NoError(t, adam.AttemptsTo(connect.ToTheNATSServerAt(natsServer.ClientURL())))

			require.NoError(t, adam.AttemptsTo(flush.TheNATSConnection()))
			assert.Equal(t, 1, natsServer.NumClients())
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the URL is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			require.ErrorIs(t, adam.AttemptsTo(connect.ToTheNATSServerAt("")), action.ErrURLIsEmpty)
		})

		t.Run("fails when the server cannot be reached", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS())
			err := adam.AttemptsTo(connect.ToTheNATSServerAt("nats://invalid:4222"))
			require.ErrorIs(t, err, ability.ErrFailedToConnect)
		})

		t.Run("fails when the actor is already connected", func(t *testing.T) {
			nc := fixture.ConnectToAnEmbeddedNATSServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
			err := adam.AttemptsTo(connect.ToTheNATSServerAt(nc.ConnectedUrl()))
			require.ErrorIs(t, err, ability.ErrAlreadyConnected)
		})

		t.Run("fails when the actor cannot communicate with NATS", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(connect.ToTheNATSServerAt("nats://localhost:4222"))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		natsServer := fixture.RunAnEmbeddedNATSServer(t)
		adam := internalfixture.NarratedActor("Adam", recorder, nil)

		require.NoError(t, adam.AttemptsTo(connect.ToTheNATSServerAt(natsServer.ClientURL())))

		assert.Equal(
			t,
			`Adam attempts to connect to the NATS server at "`+natsServer.ClientURL()+`"`,
			recorder.Messages()[0],
		)
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(
			t,
			`connect to the NATS server at "nats://localhost:4222"`,
			connect.ToTheNATSServerAt("nats://localhost:4222").String(),
		)
	})
}
