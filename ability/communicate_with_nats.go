package ability

import impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"

// CommunicateWithNATS enables the actor to communicate with a NATS server. The
// connection is given with Using, or opened by the ability itself from the URL
// given with UsingServerAt, reading "communicate with NATS using ...":
//
//	CommunicateWithNATS().Using(nc)
//	CommunicateWithNATS().UsingServerAt("nats://localhost:4222")
func CommunicateWithNATS() *impl.CommunicateWithNATSAbility {
	return impl.CommunicateWithNATS()
}
