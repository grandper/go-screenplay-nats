package been

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"
)

// Received is a resolution to tell if NATS messages have been received.
func Received() *ReceivedResolution {
	return &ReceivedResolution{}
}

// ReceivedResolution is a resolution to tell if NATS messages have been received.
type ReceivedResolution struct {
	subject string
}

// From sets the subject from which messages have been received.
func (r *ReceivedResolution) From(subject string) *ReceivedResolution {
	r.subject = subject
	return r
}

// Inspect returns nil when at least one message has been received, from the
// subject when one is set, and an error otherwise.
func (r *ReceivedResolution) Inspect(messages []*nats.Msg) error {
	if r.subject == "" {
		if len(messages) == 0 {
			return errors.New("no message has been received")
		}
		return nil
	}
	for _, msg := range messages {
		if msg != nil && msg.Subject == r.subject {
			return nil
		}
	}
	return fmt.Errorf("no message has been received from %s", r.subject)
}

// String describes the resolution's expectation.
func (r *ReceivedResolution) String() string {
	if r.subject != "" {
		return "been received from " + r.subject
	}
	return "been received"
}

// ReceivedResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution[[]*nats.Msg] = (*ReceivedResolution)(nil)
