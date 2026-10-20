package connect

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// ToTheNATSServerAt connects the actor to the NATS server at the URL. It fails
// with action.ErrURLIsEmpty when the URL is empty, with an error wrapping
// ability.ErrFailedToConnect when the server cannot be reached, and with
// ability.ErrAlreadyConnected when the actor already holds a connection.
func ToTheNATSServerAt(url string) screenplay.Performable {
	return &toTheNATSServerAtAction{url: url}
}

// toTheNATSServerAtAction connects the actor to a NATS server.
type toTheNATSServerAtAction struct {
	url string
}

// PerformAs performs the action as the provided actor.
func (a *toTheNATSServerAtAction) PerformAs(theActor *screenplay.Actor) error {
	if a.url == "" {
		return action.ErrURLIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToConnectTo(a.url)
}

// String describes the action.
func (a *toTheNATSServerAtAction) String() string {
	return "connect to the NATS server at " + verdict.Describe(a.url)
}

// toTheNATSServerAtAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*toTheNATSServerAtAction)(nil)
