package were

import (
	"errors"

	"github.com/grandper/go-screenplay/screenplay"
	"github.com/nats-io/nats.go"
)

var (
	// ErrInvalidObjectType is returned when the object is not a Connection.
	ErrInvalidObjectType = errors.New("invalid object type: expect NATS messages")
)

// Received is a matcher to tell if NATS messages were received.
func Received() *ReceivedResolution {
	return &ReceivedResolution{}
}

// ReceivedResolution is a matcher to tell if NATS messages were received.
type ReceivedResolution struct {
	subject string
}

// From sets the subject from which messages were received.
func (r *ReceivedResolution) From(subject string) *ReceivedResolution {
	r.subject = subject
	return r
}

// Resolve creates a matcher to make an assertion.
func (r *ReceivedResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		switch value := obj.(type) {
		case []*nats.Msg:
			return r.messageReceived(value), nil
		case nil:
			return false, nil
		default:
			return false, ErrInvalidObjectType
		}
	}
}

func (r *ReceivedResolution) messageReceived(messages []*nats.Msg) bool {
	if r.subject == "" {
		return len(messages) > 0
	}
	for _, msg := range messages {
		if msg.Subject == r.subject {
			return true
		}
	}
	return false
}

// String describe the resolution's expectation.
func (r *ReceivedResolution) String() string {
	return "were received"
}
