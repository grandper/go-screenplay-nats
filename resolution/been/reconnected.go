package been

import (
	"github.com/grandper/go-screenplay-nats/question"
	"github.com/grandper/go-screenplay/screenplay"
)

// Reconnected is a matcher to tell if the connection has been reconnected at some point.
func Reconnected() *ReconnectedResolution {
	return &ReconnectedResolution{}
}

// ReconnectedResolution is a matcher to tell if the connection has been reconnected at some point.
type ReconnectedResolution struct{}

// Resolve creates a matcher to make an assertion.
func (r *ReconnectedResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		switch value := obj.(type) {
		case *question.Connection:
			return value.WasReconnected(), nil
		case nil:
			return false, nil
		default:
			return false, ErrNotConnection
		}
	}
}

// String describes the resolution's expectation.
func (r *ReconnectedResolution) String() string {
	return "been reconnected"
}

// ReconnectedResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*ReconnectedResolution)(nil)
