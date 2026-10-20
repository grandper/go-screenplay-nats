package ability

import (
	"github.com/nats-io/nats.go"
)

// Message creates a new NATSMessage instance.
func Message() *NATSMessage {
	return &NATSMessage{}
}

// NATSMessage represents a message to be sent or received from NATS.
type NATSMessage struct {
	subject string
	data    []byte
	replyTo string
	headers nats.Header
}

// Subject returns the subject of the NATS message.
func (m *NATSMessage) Subject() string {
	return m.subject
}

// Data returns the data of the NATS message.
func (m *NATSMessage) Data() []byte {
	return m.data
}

// ReplyTo returns the reply-to subject of the NATS message.
func (m *NATSMessage) ReplyTo() string {
	return m.replyTo
}

// Headers returns the headers of the NATS message.
func (m *NATSMessage) Headers() nats.Header {
	return m.headers
}

// To sets the subject for the NATS message.
func (m *NATSMessage) To(subject string) *NATSMessage {
	m.subject = subject
	return m
}

// From sets the subject for the NATS message.
func (m *NATSMessage) From(subject string) *NATSMessage {
	m.subject = subject
	return m
}

// WithHeader sets a header for the NATS message.
func (m *NATSMessage) WithHeader(key, value string) *NATSMessage {
	if m.headers == nil {
		m.headers = nats.Header{}
	}
	m.headers.Add(key, value)
	return m
}

// WithHeaders sets multiple headers for the NATS message.
func (m *NATSMessage) WithHeaders(headers nats.Header) *NATSMessage {
	if m.headers == nil {
		m.headers = nats.Header{}
	}
	for key, values := range headers {
		for _, value := range values {
			m.headers.Add(key, value)
		}
	}
	return m
}

// ToBeRepliedTo sets the reply-to subject for the NATS message.
func (m *NATSMessage) ToBeRepliedTo(replyTo string) *NATSMessage {
	m.replyTo = replyTo
	return m
}

// Containing sets the data for the NATS message.
func (m *NATSMessage) Containing(data []byte) *NATSMessage {
	m.data = data
	return m
}

func (m *NATSMessage) toNATSMessage() *nats.Msg {
	return &nats.Msg{
		Subject: m.Subject(),
		Reply:   m.ReplyTo(),
		Header:  m.Headers(),
		Data:    m.Data(),
	}
}
