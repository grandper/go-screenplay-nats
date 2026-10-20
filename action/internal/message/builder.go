package message

import (
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"

	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// Builder builds the message of a publish or send action. The action holds
// one and delegates its chain methods to it, so that the message is built
// and narrated the same way by both actions.
type Builder struct {
	subject string
	data    []byte
	header  nats.Header
	reply   string
	parts   []string
	err     error // why the action cannot be performed, when it was built wrongly
}

// NewBuilder creates the builder of a message with the payload.
func NewBuilder(data []byte) *Builder {
	return &Builder{data: data}
}

// To sets the subject the message is published to; it reads "to orders".
func (b *Builder) To(subject string) {
	b.subject = subject
	b.parts = append(b.parts, "to "+verdict.Describe(subject))
}

// WithTheHeader adds a header to the message; it reads "with the header
// Key = value". An empty key is reported when the action is performed.
func (b *Builder) WithTheHeader(key, value string) {
	b.parts = append(b.parts, "with the header "+describeHeader(key, value))
	if key == "" {
		b.fail(fmt.Errorf("%w: a key cannot be empty", action.ErrInvalidHeaders))

		return
	}
	b.addHeader(key, value)
}

// WithTheHeaders adds headers to the message, written as key, value, key,
// value; it reads "with the headers A = 1, B = 2", in the order they were
// written. No header, or a key without its value, is reported when the
// action is performed.
func (b *Builder) WithTheHeaders(args ...string) {
	switch {
	case len(args) == 0:
		b.parts = append(b.parts, "with the headers")
		b.fail(fmt.Errorf("%w: WithTheHeaders should receive at least one key and its value", action.ErrInvalidHeaders))

		return
	case len(args)%2 != 0:
		b.parts = append(b.parts, "with the headers "+strings.Join(args, ", "))
		b.fail(fmt.Errorf("%w: WithTheHeaders should receive a list of key-value pairs", action.ErrInvalidHeaders))

		return
	}

	described := make([]string, 0, len(args)/2) //nolint:mnd // two arguments make one header
	for i := 0; i < len(args); i += 2 {
		if args[i] == "" {
			b.fail(fmt.Errorf("%w: a key cannot be empty", action.ErrInvalidHeaders))
		}
		b.addHeader(args[i], args[i+1])
		described = append(described, describeHeader(args[i], args[i+1]))
	}
	b.parts = append(b.parts, "with the headers "+strings.Join(described, ", "))
}

// ExpectingAReplyOn sets the subject the message is to be replied on; it
// reads "expecting a reply on inbox".
func (b *Builder) ExpectingAReplyOn(reply string) {
	b.reply = reply
	b.parts = append(b.parts, "expecting a reply on "+verdict.Describe(reply))
}

func (b *Builder) addHeader(key, value string) {
	if b.header == nil {
		b.header = nats.Header{}
	}
	b.header.Add(key, value)
}

// fail keeps the first mistake made while building the message.
func (b *Builder) fail(err error) {
	if b.err == nil {
		b.err = err
	}
}

func describeHeader(key, value string) string {
	return fmt.Sprintf("%s = %s", verdict.Describe(key), verdict.Describe(value))
}

// Err returns why the message cannot be sent: a mistake in the headers, or
// a missing subject.
func (b *Builder) Err() error {
	if b.err != nil {
		return b.err
	}
	if b.subject == "" {
		return action.ErrSubjectIsEmpty
	}

	return nil
}

// Msg returns the message to publish or to send.
func (b *Builder) Msg() *nats.Msg {
	return &nats.Msg{
		Subject: b.subject,
		Reply:   b.reply,
		Header:  b.header,
		Data:    b.data,
	}
}

// Describe returns the chain as it was written, after the name of the
// action: `to "orders" with the header "Key" = "value"`. It is empty when
// nothing was written.
func (b *Builder) Describe() string {
	return strings.Join(b.parts, " ")
}
