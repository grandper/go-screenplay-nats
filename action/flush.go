package action

import (
	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// FlushAction is an action that forces a flush on the NATS connection.
type FlushAction struct{}

// Flush forces a flush on the NATS connection.
func Flush() *FlushAction {
	return &FlushAction{}
}

// PerformAs performs the task or the action as the provided actor.
func (a *FlushAction) PerformAs(theActor *screenplay.Actor) error {
	communicateWithNATS, err := screenplay.UseAbilityTo[*ability.CommunicateWithNATSAbility]().Of(theActor)
	if err != nil {
		return err
	}
	return communicateWithNATS.Flush()
}

// String returns a string representation of the action.
func (a *FlushAction) String() string {
	return "flush the NATS connection"
}

// FlushAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*FlushAction)(nil)
