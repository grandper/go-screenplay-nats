package action

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// FlushAction is an action that forces a flush on the NATS connection.
type FlushAction struct{}

// Flush forces a flush on the NATS connection.
func Flush() *FlushAction {
	return &FlushAction{}
}

// PerformAs performs the task or the action as the provided actor.
func (a *FlushAction) PerformAs(theActor *screenplay.Actor) error {
	return impl.TheAbilityOf(theActor).ToFlush()
}

// String returns a string representation of the action.
func (a *FlushAction) String() string {
	return "flush the NATS connection"
}

// FlushAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*FlushAction)(nil)
