package ability

import impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"

// CommunicateWithNATS enables the actor to communicate with a NATS server. The
// connection is given with UsingTheConnection, or opened by the ability itself
// from the URL given with UsingTheServerAt, reading "communicate with NATS
// using the connection nc":
//
//	CommunicateWithNATS().UsingTheConnection(nc)
//	CommunicateWithNATS().UsingTheServerAt("nats://localhost:4222")
func CommunicateWithNATS() *impl.CommunicateWithNATSAbility {
	return impl.CommunicateWithNATS()
}
