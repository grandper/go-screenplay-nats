package been

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
)

// Disconnected is a resolution to tell if the connection has been disconnected at some point.
func Disconnected() *DisconnectedResolution {
	return &DisconnectedResolution{}
}

// DisconnectedResolution is a resolution to tell if the connection has been disconnected at some point.
type DisconnectedResolution struct{}

// Inspect returns nil when the connection has been disconnected, and an error otherwise.
func (r *DisconnectedResolution) Inspect(connection *ability.Connection) error {
	if connection == nil {
		return errors.New("the connection is nil")
	}
	if !connection.WasDisconnected() {
		return errors.New("it has not been disconnected")
	}
	return nil
}

// String describes the resolution's expectation.
func (r *DisconnectedResolution) String() string {
	return "been disconnected"
}

// DisconnectedResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution[*ability.Connection] = (*DisconnectedResolution)(nil)
