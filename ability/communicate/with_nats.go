package communicate

import impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"

// WithNATS enables the actor to communicate with a NATS server. The connection
// is given with UsingTheConnection, or opened by the ability itself from the
// URL given with UsingTheServerAt, reading "communicate with NATS using the
// connection nc":
//
//	communicate.WithNATS().UsingTheConnection(nc)
//	communicate.WithNATS().UsingTheServerAt("nats://localhost:4222")
func WithNATS() *impl.CommunicateWithNATSAbility {
	return impl.CommunicateWithNATS()
}
