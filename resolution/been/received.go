package been

import (
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"

	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// Received creates a resolution to tell if at least one message has been
// received; it reads "been received".
func Received() screenplay.Resolution[[]*nats.Msg] {
	return resolution.FromFunc("been received", func(messages []*nats.Msg) error {
		if len(messages) == 0 {
			return errors.New("no message has been received")
		}

		return nil
	})
}

// ReceivedFrom creates a resolution to tell if at least one message has been
// received from the subject; it reads "been received from orders". The
// subject is the one the message was published to, so a message received
// through a wildcard subscription matches its own subject.
func ReceivedFrom(subject string) screenplay.Resolution[[]*nats.Msg] {
	return resolution.FromFunc(
		"been received from "+verdict.Describe(subject),
		func(messages []*nats.Msg) error {
			for _, msg := range messages {
				if msg != nil && msg.Subject == subject {
					return nil
				}
			}

			return fmt.Errorf("no message has been received from %s", verdict.Describe(subject))
		},
	)
}
