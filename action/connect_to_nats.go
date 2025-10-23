package action

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrURLIsEmpty is returned when the URL for the NATS connection is empty.
	ErrURLIsEmpty = errors.New("the URL for the NATS connection cannot be empty")

	// ErrFailedToConnect is returned when the connection to the NATS server fails.
	ErrFailedToConnect = errors.New("failed to connect to the NATS server")
)

type ConnectToNATSAction struct {
	url string
}

// ConnectToNATS creates an action to connect to a NATS server.
func ConnectToNATS() *ConnectToNATSAction {
	return &ConnectToNATSAction{}
}

// At sets the URL of the NATS server to connect to.
func (na *ConnectToNATSAction) At(url string) *ConnectToNATSAction {
	na.url = url
	return na
}

// PerformAs performs the connection to the NATS server as the provided actor.
func (na *ConnectToNATSAction) PerformAs(theActor *screenplay.Actor) error {
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	if na.url == "" {
		return ErrURLIsEmpty
	}
	if communicateWithNATS.UsingServerAt(na.url) == nil {
		return fmt.Errorf("%w at %s", ErrFailedToConnect, na.url)
	}
	return nil
}

// String returns a string representation of the action.
func (na *ConnectToNATSAction) String() string {
	return "Connect to NATS server at " + na.url
}

// ConnectToNATSAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ConnectToNATSAction)(nil)
