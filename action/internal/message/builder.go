package message

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nats-io/nats.go"

	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// Builder builds the message of an action of type T through chain methods
// that each return the action, so that the chain keeps flowing. The action
// embeds the builder and gives it a pointer to itself.
type Builder[T any] struct {
	action  *T
	subject string
	data    []byte
	header  nats.Header
	reply   string
	parts   []string
}

// NewBuilder creates the builder of the message of the action.
func NewBuilder[T any](action *T) *Builder[T] {
	return &Builder[T]{action: action}
}

// To sets the subject the message is published to; it reads "to orders".
func (b *Builder[T]) To(subject string) *T {
	b.subject = subject
	b.parts = append(b.parts, "to "+verdict.Describe(subject))

	return b.action
}

// Containing sets the payload of the message; it reads "containing Hello".
func (b *Builder[T]) Containing(data []byte) *T {
	b.data = data
	b.parts = append(b.parts, "containing "+verdict.Describe(data))

	return b.action
}

// WithTheHeader adds a header to the message; it reads "with the header
// Key = value".
func (b *Builder[T]) WithTheHeader(key, value string) *T {
	b.addHeader(key, value)
	b.parts = append(b.parts, "with the header "+describeHeader(key, value))

	return b.action
}

// WithTheHeaders adds headers to the message; it reads "with the headers
// A = 1, B = 2", the headers sorted by key then by value.
func (b *Builder[T]) WithTheHeaders(headers nats.Header) *T {
	described := make([]string, 0, len(headers))
	for _, key := range slices.Sorted(maps(headers)) {
		for _, value := range headers[key] {
			b.addHeader(key, value)
			described = append(described, describeHeader(key, value))
		}
	}
	b.parts = append(b.parts, "with the headers "+strings.Join(described, ", "))

	return b.action
}

// maps returns the keys of the headers.
func maps(headers nats.Header) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for key := range headers {
			if !yield(key) {
				return
			}
		}
	}
}

// ToBeRepliedTo sets the subject the message is to be replied to; it reads
// "to be replied to inbox".
func (b *Builder[T]) ToBeRepliedTo(reply string) *T {
	b.reply = reply
	b.parts = append(b.parts, "to be replied to "+verdict.Describe(reply))

	return b.action
}

func (b *Builder[T]) addHeader(key, value string) {
	if b.header == nil {
		b.header = nats.Header{}
	}
	b.header.Add(key, value)
}

func describeHeader(key, value string) string {
	return fmt.Sprintf("%s = %s", verdict.Describe(key), verdict.Describe(value))
}

// Subject returns the subject of the message.
func (b *Builder[T]) Subject() string {
	return b.subject
}

// Msg returns the message to publish or to send.
func (b *Builder[T]) Msg() *nats.Msg {
	return &nats.Msg{
		Subject: b.subject,
		Reply:   b.reply,
		Header:  b.header,
		Data:    b.data,
	}
}

// Describe returns the chain as it was written, after the name of the
// action: `to "orders" containing "Hello"`. It is empty when nothing was
// written.
func (b *Builder[T]) Describe() string {
	return strings.Join(b.parts, " ")
}
