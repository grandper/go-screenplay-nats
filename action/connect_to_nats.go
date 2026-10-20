package action

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// ConnectToNATSAction is an action to connect to a NATS server.
type ConnectToNATSAction struct {
	url string
}

// ConnectToNATS creates an action to connect to a NATS server.
func ConnectToNATS() *ConnectToNATSAction {
	return &ConnectToNATSAction{}
}

// At sets the URL of the NATS server to connect to.
func (a *ConnectToNATSAction) At(url string) *ConnectToNATSAction {
	a.url = url
	return a
}

// PerformAs performs the task or the action as the provided actor.
func (a *ConnectToNATSAction) PerformAs(theActor *screenplay.Actor) error {
	if a.url == "" {
		return ErrURLIsEmpty
	}
	return impl.TheAbilityOf(theActor).ToConnectTo(a.url)
}

// String returns a string representation of the action.
func (a *ConnectToNATSAction) String() string {
	return "Connect to NATS server at " + a.url
}

// ConnectToNATSAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ConnectToNATSAction)(nil)
