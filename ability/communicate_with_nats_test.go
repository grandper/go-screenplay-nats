package ability_test

import (
	"context"
	"testing"
	"time"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay-nats/fixture"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunicateWithNATS(t *testing.T) {
	t.Run("should be created using a connection", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			a := ability.CommunicateWithNATS().Using(nc)
			require.NotNil(t, a)
			require.NoError(t, a.Forget())
		})

		t.Run("cannot be created without NATS connection", func(t *testing.T) {
			assert.Nil(t, ability.CommunicateWithNATS())
			assert.Nil(t, ability.CommunicateWithNATS().Using(nil))
		})

		t.Run("does not replace an already-initialized ability", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			a := ability.CommunicateWithNATS().Using(nc)
			require.NotNil(t, a)
			require.NoError(t, a.SubscribeTo("test.subject", 0))
			assert.Equal(t, 1, nc.NumSubscriptions())

			same := a.Using(nc)
			assert.Same(t, a, same)
			assert.Equal(t, 1, nc.NumSubscriptions())
			require.NoError(t, a.Forget())
		})
	})

	t.Run("can be created from a NATS server url", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			require.NotNil(t, ability.CommunicateWithNATS().UsingServerAt(nc.ConnectedUrl()))
		})

		t.Run("cannot be created without empty url", func(t *testing.T) {
			require.Nil(t, ability.CommunicateWithNATS().UsingServerAt(""))
		})

		t.Run("cannot be created with invalid url", func(t *testing.T) {
			require.Nil(t, ability.CommunicateWithNATS().UsingServerAt("foobar"))
		})

		t.Run("does not replace an already-initialized ability", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			a := ability.CommunicateWithNATS().Using(nc)
			require.NotNil(t, a)
			require.NoError(t, a.SubscribeTo("test.subject", 0))
			assert.Equal(t, 1, nc.NumSubscriptions())

			same := a.UsingServerAt(nc.ConnectedUrl())
			assert.Same(t, a, same)
			assert.Equal(t, 1, nc.NumSubscriptions())
			require.NoError(t, a.Forget())
		})
	})

	t.Run("provide connection details when it is", func(t *testing.T) {
		t.Run("closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.False(t, communicateWithNATS.ConnectionWasClosed())

			nc.Close()
			assertEventually(t, func() bool {
				return nc.IsClosed()
			})
			assert.True(t, communicateWithNATS.ConnectionWasClosed())
		})

		t.Run("disconnected", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.False(t, communicateWithNATS.ConnectionWasDisconnected())

			assert.NoError(t, nc.ForceReconnect())
			assertEventually(t, func() bool {
				return nc.IsConnected()
			})
			assert.True(t, communicateWithNATS.ConnectionWasDisconnected())
		})

		t.Run("reconnected", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.False(t, communicateWithNATS.ConnectionWasReconnected())

			assert.NoError(t, nc.ForceReconnect())
			assertEventually(t, func() bool {
				return nc.IsConnected()
			})
			assert.True(t, communicateWithNATS.ConnectionWasReconnected())
		})

		t.Run("producing an async error", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()
			sub, err := nc.SubscribeSync("test")
			require.NoError(t, err)
			require.NoError(t, sub.SetPendingLimits(1, 1))

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.False(t, communicateWithNATS.ConnectionAsyncErrorOccurred())

			for i := 0; i < 10; i++ {
				require.NoError(t, communicateWithNATS.Publish(ability.Message().To("test").Containing([]byte("msg"))))
			}
			require.NoError(t, nc.Flush())
			assertEventually(t, func() bool { return communicateWithNATS.ConnectionAsyncErrorOccurred() })
		})
	})

	subject := "test.subject"
	queue := "test.queue"
	message := []byte("hello world")

	t.Run("publishes messages", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			require.NoError(t, communicateWithNATS.Publish(ability.Message().To(subject).Containing(message)))

			msg, err := sub.NextMsg(testTimeout)
			assert.NoError(t, err)
			assert.Equal(t, message, msg.Data)
			assert.Equal(t, subject, msg.Subject)
		})

		t.Run("fails when the message is nil", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			err := communicateWithNATS.Publish(nil)
			assert.ErrorIs(t, err, ability.ErrMessageIsNil)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			err := communicateWithNATS.Publish(ability.Message().To(subject).Containing(message))
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			err := communicateWithNATS.Publish(ability.Message().Containing(message))
			assert.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
		})
	})

	t.Run("responds to requests", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			sub, err := nc.Subscribe("test.subject", func(msg *nats.Msg) {
				assert.NoError(t, communicateWithNATS.Respond(msg, "response"))
			})
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()

			msg, err := nc.Request("test.subject", message, time.Second)
			assert.NoError(t, err)
			assert.Equal(t, []byte("response"), msg.Data)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			msg := &nats.Msg{Subject: subject}
			err := communicateWithNATS.Respond(msg, "response")
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the message is nil", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			err := communicateWithNATS.Respond(nil, "response")
			assert.ErrorIs(t, err, ability.ErrMessageIsNil)
		})
	})

	t.Run("flushes the connection", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			require.NoError(t, communicateWithNATS.Publish(ability.Message().To(subject).Containing(message)))
			assert.NoError(t, communicateWithNATS.Flush())
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			err := communicateWithNATS.Flush()
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})
	})

	t.Run("sends requests", func(t *testing.T) {
		const requestTimeout = time.Second
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				assert.NoError(t, msg.Respond([]byte("ok")))
			})
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			response, err := communicateWithNATS.Request(ability.Message().To(subject).Containing(message), requestTimeout)
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, response.Data, []byte("ok"))
		})

		t.Run("fails when the request is nil", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			response, err := communicateWithNATS.Request(nil, requestTimeout)
			assert.ErrorIs(t, err, ability.ErrMessageIsNil)
			assert.Nil(t, response)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			response, err := communicateWithNATS.Request(ability.Message().To(subject).Containing(message), requestTimeout)
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
			assert.Nil(t, response)
		})

		t.Run("fails when subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			response, err := communicateWithNATS.Request(ability.Message().Containing(message), requestTimeout)
			assert.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
			assert.Nil(t, response)
		})

		t.Run("fails when no responder is set", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			response, err := communicateWithNATS.Request(ability.Message().To(subject).Containing(message), requestTimeout)
			assert.ErrorIs(t, err, nats.ErrNoResponders)
			assert.Nil(t, response)
		})

		t.Run("fails when the request takes too long", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				// No response
			})
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, sub.Unsubscribe())
			}()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			_, err = communicateWithNATS.Request(ability.Message().To(subject).Containing(message), requestTimeout)
			assert.ErrorIs(t, err, nats.ErrTimeout)
		})
	})

	t.Run("subscribes to a subject asynchronously", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			assert.Equal(t, 0, nc.NumSubscriptions())

			err := communicateWithNATS.SubscribeAsyncTo(subject, 0)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())
		})

		t.Run("stores received messages", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			require.NoError(t, communicateWithNATS.SubscribeAsyncTo(subject, 0))
			require.NoError(t, communicateWithNATS.Publish(ability.Message().To(subject).Containing(message)))

			assert.Eventually(t, func() bool {
				msg, err := communicateWithNATS.LastMessageForSubject(subject)
				return err == nil && string(msg.Data) == string(message)
			}, testTimeout, tickInterval)
		})

		t.Run("successfully with auto unsubscribe", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.SubscribeAsyncTo(subject, 5)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())

			for i := 0; i < 5; i++ {
				err = communicateWithNATS.Publish(ability.Message().To(subject).Containing([]byte("msg")))
				assert.NoError(t, err)
				assert.Equal(t, 1, nc.NumSubscriptions())
			}

			err = communicateWithNATS.Publish(ability.Message().To(subject).Containing([]byte("msg")))
			assert.NoError(t, err)
			assert.Eventually(t, func() bool {
				return nc.NumSubscriptions() == 0
			}, testTimeout, tickInterval)
		})

		t.Run("fails if the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeAsyncTo("", 0)
			assert.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeAsyncTo(subject, 0)
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the subscription already exists", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeAsyncTo(subject, 0)
			assert.NoError(t, err)

			err = communicateWithNATS.SubscribeAsyncTo(subject, 0)
			assert.ErrorIs(t, err, ability.ErrSubscriptionAlreadyExists)
		})

		t.Run("fails when the connection is closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			nc.Close()

			err := communicateWithNATS.SubscribeAsyncTo(subject, 0)
			assert.Error(t, err)
		})
	})

	t.Run("subscribes to a subject synchronously", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			assert.Equal(t, 0, nc.NumSubscriptions())

			err := communicateWithNATS.SubscribeTo(subject, 0)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())
		})

		t.Run("successfully with auto unsubscribe", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.SubscribeTo(subject, 5)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())

			for i := 0; i < 5; i++ {
				err = communicateWithNATS.Publish(ability.Message().To(subject).Containing([]byte("msg")))
				assert.NoError(t, err)
				assert.Equal(t, 1, nc.NumSubscriptions())

			}
			err = communicateWithNATS.ReceiveAllMessages(context.Background(), 200*time.Millisecond)
			require.ErrorIs(t, err, nats.ErrMaxMessages)
			assert.Eventually(t, func() bool {
				return nc.NumSubscriptions() == 0
			}, testTimeout, tickInterval)
		})

		t.Run("fails if the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeTo("", 0)
			assert.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeTo(subject, 0)
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the subscription already exists", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeTo(subject, 0)
			assert.NoError(t, err)

			err = communicateWithNATS.SubscribeTo(subject, 0)
			assert.ErrorIs(t, err, ability.ErrSubscriptionAlreadyExists)
		})

		t.Run("fails when the connection is closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			nc.Close()

			err := communicateWithNATS.SubscribeTo(subject, 0)
			assert.Error(t, err)
		})
	})

	t.Run("subscribes to a queue asynchronously", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			assert.Equal(t, 0, nc.NumSubscriptions())

			err := communicateWithNATS.SubscribeToQueueAsync(subject, queue, 0)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())
		})

		t.Run("successfully with auto unsubscribe", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.SubscribeToQueueAsync(subject, queue, 5)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())

			for i := 0; i < 5; i++ {
				err = communicateWithNATS.Publish(ability.Message().To(subject).Containing([]byte("msg")))
				assert.NoError(t, err)
				assert.Equal(t, 1, nc.NumSubscriptions())
			}

			err = communicateWithNATS.Publish(ability.Message().To(subject).Containing([]byte("msg")))
			assert.NoError(t, err)
			assert.Eventually(t, func() bool {
				return nc.NumSubscriptions() == 0
			}, testTimeout, tickInterval)
		})

		t.Run("fails if the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueueAsync("", queue, 0)
			assert.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
		})

		t.Run("fails if the queue is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueueAsync(subject, "", 0)
			assert.ErrorIs(t, err, ability.ErrQueueIsEmpty)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueueAsync(subject, queue, 0)
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the subscription already exists", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueueAsync(subject, queue, 0)
			assert.NoError(t, err)

			err = communicateWithNATS.SubscribeToQueueAsync(subject, queue, 0)
			assert.ErrorIs(t, err, ability.ErrSubscriptionAlreadyExists)
		})

		t.Run("fails when the connection is closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			nc.Close()

			err := communicateWithNATS.SubscribeToQueueAsync(subject, queue, 0)
			assert.Error(t, err)
		})
	})

	t.Run("subscribes to a queue synchronously", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			assert.Equal(t, 0, nc.NumSubscriptions())

			err := communicateWithNATS.SubscribeToQueue(subject, queue, 0)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())
		})

		t.Run("successfully with auto unsubscribe", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.SubscribeToQueue(subject, queue, 5)
			assert.NoError(t, err)

			assert.Equal(t, 1, nc.NumSubscriptions())

			for i := 0; i < 5; i++ {
				err = communicateWithNATS.Publish(ability.Message().To(subject).Containing([]byte("msg")))
				assert.NoError(t, err)
				assert.Equal(t, 1, nc.NumSubscriptions())

			}
			err = communicateWithNATS.ReceiveAllMessages(context.Background(), 200*time.Millisecond)
			require.ErrorIs(t, err, nats.ErrMaxMessages)
			assert.Eventually(t, func() bool {
				return nc.NumSubscriptions() == 0
			}, testTimeout, tickInterval)
		})

		t.Run("fails if the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueue("", queue, 0)
			assert.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
		})

		t.Run("fails if the queue is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueue(subject, "", 0)
			assert.ErrorIs(t, err, ability.ErrQueueIsEmpty)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueue(subject, queue, 0)
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the subscription already exists", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			err := communicateWithNATS.SubscribeToQueue(subject, queue, 0)
			assert.NoError(t, err)

			err = communicateWithNATS.SubscribeToQueue(subject, queue, 0)
			assert.ErrorIs(t, err, ability.ErrSubscriptionAlreadyExists)
		})

		t.Run("fails when the connection is closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			nc.Close()

			err := communicateWithNATS.SubscribeToQueue(subject, queue, 0)
			assert.Error(t, err)
		})
	})

	t.Run("unsubscribes from subjects", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.SubscribeAsyncTo(subject, 0)
			assert.NoError(t, err)
			assert.Equal(t, 1, nc.NumSubscriptions())

			err = communicateWithNATS.UnsubscribeFrom(subject)
			assert.NoError(t, err)
			assert.Equal(t, 0, nc.NumSubscriptions())
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.UnsubscribeFrom("")
			assert.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.UnsubscribeFrom(subject)
			assert.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when there is no subscription for the subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.UnsubscribeFrom(subject)
			assert.ErrorIs(t, err, ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the connection was closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()

			err := communicateWithNATS.SubscribeAsyncTo(subject, 0)
			assert.NoError(t, err)
			assert.Equal(t, 1, nc.NumSubscriptions())

			nc.Close()

			err = communicateWithNATS.UnsubscribeFrom(subject)
			assert.Error(t, err)
		})
	})

	t.Run("receives the next message for a subject", func(t *testing.T) {
		ctx := context.Background()

		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))

			assert.NoError(t, nc.Publish(subject, []byte("hello world")))

			assert.NoError(t, communicateWithNATS.ReceiveNextMessageFor(ctx, subject))
			msg, err := communicateWithNATS.LastMessageForSubject(subject)
			assert.NoError(t, err)
			assert.Equal(t, []byte("hello world"), msg.Data)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.ErrorIs(t, communicateWithNATS.ReceiveNextMessageFor(ctx, ""), ability.ErrSubjectIsEmpty)
		})

		t.Run("fails then the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			assert.ErrorIs(t, communicateWithNATS.ReceiveNextMessageFor(ctx, subject), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when there is no subscription for the subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.ErrorIs(t, communicateWithNATS.ReceiveNextMessageFor(ctx, subject), ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the connection was closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))
			assert.NoError(t, nc.Publish(subject, []byte("hello world")))
			nc.Close()

			assert.Error(t, communicateWithNATS.ReceiveNextMessageFor(ctx, subject))
		})
	})

	t.Run("receives all the messages for a subject", func(t *testing.T) {
		ctx := context.Background()
		timeout := 20 * time.Millisecond

		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))

			assert.NoError(t, nc.Publish(subject, []byte("hello world")))
			assert.NoError(t, nc.Publish(subject, []byte("bye world")))

			assert.NoError(t, communicateWithNATS.ReceiveAllMessagesFor(ctx, subject, timeout))
			msgs, err := communicateWithNATS.AllMessagesForSubject(subject)
			assert.NoError(t, err)
			assert.Len(t, msgs, 2)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.ErrorIs(t, communicateWithNATS.ReceiveAllMessagesFor(ctx, "", timeout), ability.ErrSubjectIsEmpty)
		})

		t.Run("fails then the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			assert.ErrorIs(t, communicateWithNATS.ReceiveAllMessagesFor(ctx, subject, timeout), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when there is no subscription for the subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.ErrorIs(t, communicateWithNATS.ReceiveAllMessagesFor(ctx, subject, timeout), ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the connection was closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))
			assert.NoError(t, nc.Publish(subject, []byte("hello world")))
			nc.Close()

			assert.Error(t, communicateWithNATS.ReceiveAllMessagesFor(ctx, subject, timeout))
		})
	})

	t.Run("receives all the messages", func(t *testing.T) {
		ctx := context.Background()
		timeout := 20 * time.Millisecond

		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))

			assert.NoError(t, nc.Publish(subject, []byte("hello world")))
			assert.NoError(t, nc.Publish(subject, []byte("bye world")))

			assert.NoError(t, communicateWithNATS.ReceiveAllMessages(ctx, timeout))
			msgs, err := communicateWithNATS.AllMessagesForSubject(subject)
			assert.NoError(t, err)
			assert.Len(t, msgs, 2)
		})

		t.Run("fails then the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			assert.ErrorIs(t, communicateWithNATS.ReceiveAllMessages(ctx, timeout), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the connection was closed", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))
			assert.NoError(t, nc.Publish(subject, []byte("hello world")))
			nc.Close()

			assert.Error(t, communicateWithNATS.ReceiveAllMessages(ctx, timeout))
		})
	})

	t.Run("provides the last message for a subject", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			require.NoError(t, communicateWithNATS.Publish(ability.Message().To(subject).Containing(message)))
			require.NoError(t, communicateWithNATS.ReceiveNextMessageFor(context.Background(), subject))

			msg, err := communicateWithNATS.LastMessageForSubject(subject)
			require.NoError(t, err)
			assert.Equal(t, message, msg.Data)

			msg, err = communicateWithNATS.LastMessageForSubject("unknown")
			require.ErrorIs(t, err, ability.ErrNoMessageForSubject)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			msg, err := communicateWithNATS.LastMessageForSubject("")
			require.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
			assert.Nil(t, msg)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			msg, err := communicateWithNATS.LastMessageForSubject(subject)
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
			assert.Nil(t, msg)
		})

		t.Run("fails when there is no message for the subject", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			msg, err := communicateWithNATS.LastMessageForSubject(subject)
			require.ErrorIs(t, err, ability.ErrNoMessageForSubject)
			assert.Nil(t, msg)
		})
	})

	t.Run("provides the messages for a subject", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			require.NoError(t, communicateWithNATS.Publish(ability.Message().To(subject).Containing(message)))
			require.NoError(t, communicateWithNATS.ReceiveNextMessageFor(context.Background(), subject))

			messages, err := communicateWithNATS.AllMessagesForSubject(subject)
			require.NoError(t, err)
			assert.Len(t, messages, 1)

			messages2, err := communicateWithNATS.AllMessagesForSubject("unknown")
			require.NoError(t, err)
			assert.Nil(t, messages2)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			messages, err := communicateWithNATS.AllMessagesForSubject(subject)
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
			assert.Nil(t, messages)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			messages, err := communicateWithNATS.AllMessagesForSubject("")
			require.ErrorIs(t, err, ability.ErrSubjectIsEmpty)
			assert.Nil(t, messages)
		})
	})

	t.Run("forgets its resources", func(t *testing.T) {
		t.Run("closes the connection", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			require.NoError(t, communicateWithNATS.Forget())
			assert.True(t, nc.IsClosed())
		})

		t.Run("unsubscribes all subscriptions", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			require.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))
			assert.Equal(t, 1, nc.NumSubscriptions())
			require.NoError(t, communicateWithNATS.Forget())
			assert.Equal(t, 0, nc.NumSubscriptions())
		})

		t.Run("does nothing when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			require.NoError(t, communicateWithNATS.Forget())
		})
	})

	t.Run("provides all the messages", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			defer nc.Close()

			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			assert.NoError(t, communicateWithNATS.SubscribeTo(subject, 0))
			defer func() {
				assert.NoError(t, communicateWithNATS.Forget())
			}()
			require.NoError(t, communicateWithNATS.Publish(ability.Message().To(subject).Containing(message)))
			require.NoError(t, communicateWithNATS.ReceiveNextMessageFor(context.Background(), subject))

			messages, err := communicateWithNATS.AllMessages()
			require.NoError(t, err)
			assert.Len(t, messages, 1)
		})

		t.Run("fails when the connection is nil", func(t *testing.T) {
			communicateWithNATS := ability.CommunicateWithNATS().Using(nil)
			msg, err := communicateWithNATS.AllMessages()
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
			assert.Nil(t, msg)
		})
	})
}

const testTimeout = 2 * time.Second
const tickInterval = 100 * time.Millisecond

func assertEventually(t *testing.T, condition func() bool) bool {
	t.Helper()
	return assert.Eventually(t, condition, testTimeout, tickInterval)
}
