package action

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
)

// ErrURLIsEmpty is returned when the URL for the NATS connection is empty.
var ErrURLIsEmpty = errors.New("the URL for the NATS connection cannot be empty")

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
	communicateWithNATS := screenplay.UseAnAbilityOf(theActor).That(ability.CommunicateWithNATS())
	if communicateWithNATS == nil {
		return screenplay.ErrActorMissingAbility
	}
	if a.url == "" {
		return ErrURLIsEmpty
	}
	return communicateWithNATS.ConnectTo(a.url)
}

// String returns a string representation of the action.
func (a *ConnectToNATSAction) String() string {
	return "Connect to NATS server at " + a.url
}

// ConnectToNATSAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ConnectToNATSAction)(nil)
