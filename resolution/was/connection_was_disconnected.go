package was

import (
	"github.com/grandper/go-screenplay-nats/question"
	"github.com/grandper/go-screenplay/screenplay"
)

// Disconnected is a matcher to tell if the connection was disconnected at some point.
func Disconnected() *DisconnectedResolution {
	return &DisconnectedResolution{}
}

// DisconnectedResolution is a matcher to tell if the connection was disconnected at some point.
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
			return false, ErrInvalidObjectType
		}
	}
}

// String describe the resolution's expectation.
func (r *DisconnectedResolution) String() string {
	return "was disconnected"
}
