package ability

// Connection is a snapshot of what happened to the NATS connection of the
// actor, taken when a question is answered. It does not change afterwards:
// ask the question again to observe a later state.
type Connection struct {
	wasClosed       bool
	wasDisconnected bool
	wasReconnected  bool
}

// NewConnection creates a snapshot of a connection that was closed,
// disconnected or reconnected, or none of these.
func NewConnection(wasClosed, wasDisconnected, wasReconnected bool) *Connection {
	return &Connection{
		wasClosed:       wasClosed,
		wasDisconnected: wasDisconnected,
		wasReconnected:  wasReconnected,
	}
}

// WasClosed tells if the connection was closed at some point.
func (c *Connection) WasClosed() bool {
	return c.wasClosed
}

// WasDisconnected tells if the connection was disconnected at some point.
func (c *Connection) WasDisconnected() bool {
	return c.wasDisconnected
}

// WasReconnected tells if the connection was reconnected at some point.
func (c *Connection) WasReconnected() bool {
	return c.wasReconnected
}
