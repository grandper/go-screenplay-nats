package been

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/question"
)

// Closed is a resolution to tell if the connection has been closed at some point.
func Closed() *ClosedResolution {
	return &ClosedResolution{}
}

// ClosedResolution is a resolution to tell if the connection has been closed at some point.
type ClosedResolution struct{}

// Inspect returns nil when the connection has been closed, and an error otherwise.
func (r *ClosedResolution) Inspect(connection *question.Connection) error {
	if connection == nil {
		return errors.New("the connection is nil")
	}
	if !connection.WasClosed() {
		return errors.New("it has not been closed")
	}
	return nil
}

// String describes the resolution's expectation.
func (r *ClosedResolution) String() string {
	return "been closed"
}

// ClosedResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution[*question.Connection] = (*ClosedResolution)(nil)
