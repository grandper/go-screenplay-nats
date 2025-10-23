package ability

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	t.Run("creates a new NATS message", func(t *testing.T) {
		msg := Message()
		require.NotNil(t, msg)

		assert.Equal(t, "", msg.Subject())
		assert.Equal(t, "", msg.ReplyTo())
		assert.Nil(t, msg.Headers())
		assert.Nil(t, msg.Data())
	})

	t.Run("sets and gets subject", func(t *testing.T) {
		subject := "test.subject"
		msgTo := Message().To(subject)
		require.Equal(t, subject, msgTo.Subject())

		msgFrom := Message().From(subject)
		require.Equal(t, subject, msgFrom.Subject())
	})

	t.Run("sets and gets and gets reply to", func(t *testing.T) {
		replyTo := "reply.subject"
		msg := Message().ToBeRepliedTo(replyTo)
		require.Equal(t, replyTo, msg.ReplyTo())
	})

	t.Run("sets and gets a header", func(t *testing.T) {
		key := "Content-Type"
		value := "application/json"
		msg := Message().WithHeader(key, value)
		require.NotNil(t, msg.Headers())
		require.Equal(t, value, msg.Headers().Get(key))
	})

	t.Run("sets and gets several headers", func(t *testing.T) {
		msg := Message().
			WithHeaders(nats.Header{
				"Key1": []string{"Value1"},
				"Key2": []string{"Value2"},
			})
		require.NotNil(t, msg.Headers())
		require.Equal(t, "Value1", msg.Headers().Get("Key1"))
		require.Equal(t, "Value2", msg.Headers().Get("Key2"))
	})

	t.Run("sets and gets a data", func(t *testing.T) {
		data := []byte("test data")
		msg := Message().Containing(data)
		require.Equal(t, data, msg.Data())
	})
}
