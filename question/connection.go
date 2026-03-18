package question

// Connection represents a point-in-time snapshot of a NATS connection's state.
// The values are captured at the moment the NATSConnection question is answered
// and will not reflect any state changes that occur after that point.
// Do not cache a *Connection across assertions; re-ask the question each time
// to obtain a fresh snapshot.
type Connection struct {
	wasClosed       bool
	wasDisconnected bool
	wasReconnected  bool
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
