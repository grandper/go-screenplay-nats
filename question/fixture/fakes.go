package fixture

import (
	"github.com/grandper/go-screenplay-nats/v2/ability"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// CreateFakeConnection creates a fake connection nothing happened to.
func CreateFakeConnection() *ability.Connection {
	return impl.NewConnection(false, false, false)
}

// CreateFakeClosedConnection creates a fake connection that was disconnected
// and then closed.
func CreateFakeClosedConnection() *ability.Connection {
	return impl.NewConnection(true, true, false)
}

// CreateFakeDisconnectedConnection creates a fake connection that was
// disconnected and not reconnected yet.
func CreateFakeDisconnectedConnection() *ability.Connection {
	return impl.NewConnection(false, true, false)
}

// CreateFakeReconnectedConnection creates a fake connection that was
// disconnected and then reconnected.
func CreateFakeReconnectedConnection() *ability.Connection {
	return impl.NewConnection(false, true, true)
}
