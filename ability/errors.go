package ability

import impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"

var (
	// ErrNATSConnectionNotSet is returned when the ability has neither a NATS
	// connection nor the URL of a server to connect to.
	ErrNATSConnectionNotSet = impl.ErrNATSConnectionNotSet

	// ErrFailedToConnect is wrapped by the error returned when the connection to
	// the NATS server fails.
	ErrFailedToConnect = impl.ErrFailedToConnect

	// ErrAlreadyConnected is returned when the actor connects while its ability
	// already holds a connection.
	ErrAlreadyConnected = impl.ErrAlreadyConnected

	// ErrSubscriptionAlreadyExists is returned when a subscription for the subject already exists.
	ErrSubscriptionAlreadyExists = impl.ErrSubscriptionAlreadyExists

	// ErrNoSubscriptionFound is wrapped by the error returned when no subscription
	// is found for the subject.
	ErrNoSubscriptionFound = impl.ErrNoSubscriptionFound
)
