package flush

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// TheNATSConnection flushes the buffered publish operations of the NATS
// connection: the actor waits until the server has received everything it
// published or subscribed to so far.
func TheNATSConnection() screenplay.Performable {
	return &theNATSConnectionAction{}
}

// theNATSConnectionAction flushes the NATS connection.
type theNATSConnectionAction struct{}

// PerformAs performs the action as the provided actor.
func (a *theNATSConnectionAction) PerformAs(theActor *screenplay.Actor) error {
	return impl.TheAbilityOf(theActor).ToFlush()
}

// String describes the action.
func (a *theNATSConnectionAction) String() string {
	return "flush the NATS connection"
}

// theNATSConnectionAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*theNATSConnectionAction)(nil)
