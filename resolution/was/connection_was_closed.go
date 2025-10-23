package was

import (
	"errors"

	"github.com/grandper/go-screenplay-nats/question"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrInvalidObjectType is returned when the object is not a Connection.
	ErrInvalidObjectType = errors.New("invalid object type: expect a Connection")
)

// Closed is a matcher to tell if the connection was closed at some point.
func Closed() *ClosedResolution {
	return &ClosedResolution{}
}

// ClosedResolution is a matcher to tell if the connection was closed at some point.
type ClosedResolution struct{}

// Resolve creates a matcher to make an assertion.
func (r *ClosedResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		switch value := obj.(type) {
		case *question.Connection:
			return value.WasClosed(), nil
		case nil:
			return false, nil
		default:
			return false, ErrInvalidObjectType
		}
	}
}

// String describe the resolution's expectation.
func (r *ClosedResolution) String() string {
	return "was closed"
}
