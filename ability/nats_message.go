package ability

import impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"

// NATSMessage represents a message to be sent to NATS. It is built fluently:
//
//	ability.Message().To("orders").Containing([]byte("Hello")).WithHeader("key", "value")
type NATSMessage = impl.NATSMessage

// Message creates a new NATSMessage.
func Message() *impl.NATSMessage {
	return impl.Message()
}
