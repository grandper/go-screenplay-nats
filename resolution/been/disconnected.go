package been

import (
	"github.com/grandper/go-screenplay/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/question"
)

// Disconnected is a matcher to tell if the connection has been disconnected at some point.
func Disconnected() *DisconnectedResolution {
	return &DisconnectedResolution{}
}

// DisconnectedResolution is a matcher to tell if the connection has been disconnected at some point.
type DisconnectedResolution struct{}

// Resolve creates a matcher to make an assertion.
func (r *DisconnectedResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		switch value := obj.(type) {
		case *question.Connection:
			return value.WasDisconnected(), nil
		case nil:
			return false, nil
		default:
			return false, ErrNotConnection
		}
	}
}

// String describes the resolution's expectation.
func (r *DisconnectedResolution) String() string {
	return "been disconnected"
}

// DisconnectedResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*DisconnectedResolution)(nil)
