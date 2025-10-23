package question

// Connection represents a connection to a NATS server.
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
