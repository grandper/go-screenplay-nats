package ability

import impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"

// Connection is a snapshot of what happened to the NATS connection of the
// actor, read through WasClosed(), WasDisconnected() and WasReconnected(). It
// is taken when a question is answered and does not change afterwards: ask the
// question again to observe a later state.
type Connection = impl.Connection
