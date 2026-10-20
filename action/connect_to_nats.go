package action

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
)

var (
	// ErrURLIsEmpty is returned when the URL for the NATS connection is empty.
	ErrURLIsEmpty = errors.New("the URL for the NATS connection cannot be empty")

	// ErrFailedToConnect is returned when the connection to the NATS server fails.
	ErrFailedToConnect = errors.New("failed to connect to the NATS server")
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
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	if a.url == "" {
		return ErrURLIsEmpty
	}
	newAbility := communicateWithNATS.UsingServerAt(a.url)
	if newAbility == nil {
		return fmt.Errorf("%w at %s", ErrFailedToConnect, a.url)
	}
	theActor.WhoCan(newAbility)
	return nil
}

// String returns a string representation of the action.
func (a *ConnectToNATSAction) String() string {
	return "Connect to NATS server at " + a.url
}

// ConnectToNATSAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ConnectToNATSAction)(nil)
