package action

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
)

// FlushAction is an action that forces a flush on the NATS connection.
type FlushAction struct{}

// Flush forces a flush on the NATS connection.
func Flush() *FlushAction {
	return &FlushAction{}
}

// PerformAs performs the task or the action as the provided actor.
func (a *FlushAction) PerformAs(theActor *screenplay.Actor) error {
	communicateWithNATS := screenplay.UseAnAbilityOf(theActor).That(ability.CommunicateWithNATS())
	return communicateWithNATS.Flush()
}

// String returns a string representation of the action.
func (a *FlushAction) String() string {
	return "flush the NATS connection"
}

// FlushAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*FlushAction)(nil)
