package been

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
)

// Reconnected is a resolution to tell if the connection has been reconnected at some point.
func Reconnected() *ReconnectedResolution {
	return &ReconnectedResolution{}
}

// ReconnectedResolution is a resolution to tell if the connection has been reconnected at some point.
type ReconnectedResolution struct{}

// Inspect returns nil when the connection has been reconnected, and an error otherwise.
func (r *ReconnectedResolution) Inspect(connection *ability.Connection) error {
	if connection == nil {
		return errors.New("the connection is nil")
	}
	if !connection.WasReconnected() {
		return errors.New("it has not been reconnected")
	}
	return nil
}

// String describes the resolution's expectation.
func (r *ReconnectedResolution) String() string {
	return "been reconnected"
}

// ReconnectedResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution[*ability.Connection] = (*ReconnectedResolution)(nil)
