package been

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/ability"
)

// errConnectionIsNil is returned when the inspected connection is nil.
var errConnectionIsNil = errors.New("the connection is nil")

// Closed creates a resolution to tell if the connection has been closed at
// some point; it reads "been closed".
func Closed() screenplay.Resolution[*ability.Connection] {
	return connection("been closed", "it has not been closed", (*ability.Connection).WasClosed)
}

// Disconnected creates a resolution to tell if the connection has been
// disconnected at some point; it reads "been disconnected". A connection that
// has been closed or reconnected has been disconnected first.
func Disconnected() screenplay.Resolution[*ability.Connection] {
	return connection("been disconnected", "it has not been disconnected", (*ability.Connection).WasDisconnected)
}

// Reconnected creates a resolution to tell if the connection has been
// reconnected at some point; it reads "been reconnected".
func Reconnected() screenplay.Resolution[*ability.Connection] {
	return connection("been reconnected", "it has not been reconnected", (*ability.Connection).WasReconnected)
}

// connection creates a resolution satisfied when happened reports true for
// the inspected connection. A nil connection never satisfies it.
func connection(
	description, failure string,
	happened func(*ability.Connection) bool,
) screenplay.Resolution[*ability.Connection] {
	return resolution.FromFunc(description, func(connection *ability.Connection) error {
		if connection == nil {
			return errConnectionIsNil
		}
		if !happened(connection) {
			return errors.New(failure)
		}

		return nil
	})
}
