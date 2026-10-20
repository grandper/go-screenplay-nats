package ability_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/fixture"
	"github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

const (
	subject        = "test.subject"
	queue          = "test.queue"
	requestTimeout = time.Second
	testTimeout    = 2 * time.Second
	tickInterval   = 100 * time.Millisecond
)

func TestCommunicateWithNATSAbility(t *testing.T) {
	ctx := context.Background()
	message := []byte("hello world")

	t.Run("is given a connection with Using", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))
			assert.Equal(t, 1, nc.NumSubscriptions())
		})

		t.Run("has no connection until one is given", func(t *testing.T) {
			require.ErrorIs(t, abilityWithoutConnection().ToFlush(), ability.ErrNATSConnectionNotSet)
			theAbility := ability.TheAbilityOf(
				screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nil)),
			)
			require.ErrorIs(t, theAbility.ToFlush(), ability.ErrNATSConnectionNotSet)
		})

		t.Run("keeps the first connection it was given", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			other := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			theAbility := ability.TheAbilityOf(screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS))
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))

			assert.Same(t, communicateWithNATS, communicateWithNATS.Using(other))
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: "other.subject"}))
			assert.Equal(t, 2, nc.NumSubscriptions())
			assert.Equal(t, 0, other.NumSubscriptions())
		})

		t.Run("ignores a nil receiver", func(t *testing.T) {
			var communicateWithNATS *ability.CommunicateWithNATSAbility
			assert.Nil(t, communicateWithNATS.Using(fixture.CreateEmbeddedNATSServer(t)))
			assert.Nil(t, communicateWithNATS.UsingServerAt("nats://localhost:4222"))
		})
	})

	t.Run("is given the URL of a server with UsingServerAt", func(t *testing.T) {
		t.Run("connects when it is first used", func(t *testing.T) {
			natsServer := fixture.RunNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().UsingServerAt(natsServer.ClientURL())
			theAbility := ability.TheAbilityOf(screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS))
			assert.Equal(t, 0, natsServer.NumClients())

			require.NoError(t, theAbility.ToFlush())
			assertEventually(t, func() bool { return natsServer.NumClients() == 1 })
			require.NoError(t, communicateWithNATS.Forget())
		})

		t.Run("fails when it is used with an empty URL", func(t *testing.T) {
			theAbility := ability.TheAbilityOf(
				screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingServerAt("")),
			)
			require.ErrorIs(t, theAbility.ToFlush(), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when it is used with an unreachable URL", func(t *testing.T) {
			theAbility := ability.TheAbilityOf(
				screenplay.ActorNamed("Adam").
					WhoCan(ability.CommunicateWithNATS().UsingServerAt("nats://invalid:4222")),
			)
			require.ErrorIs(t, theAbility.ToFlush(), ability.ErrFailedToConnect)
		})

		t.Run("keeps the connection it was given", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			other := fixture.RunNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			theAbility := ability.TheAbilityOf(screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS))

			assert.Same(t, communicateWithNATS, communicateWithNATS.UsingServerAt(other.ClientURL()))
			require.NoError(t, theAbility.ToFlush())
			assert.Equal(t, 0, other.NumClients())
		})
	})

	t.Run("connects to a server", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			natsServer := fixture.RunNATSServer(t)
			theAbility := abilityWithoutConnection()
			require.NoError(t, theAbility.ToConnectTo(natsServer.ClientURL()))
			assertEventually(t, func() bool { return natsServer.NumClients() == 1 })
			require.NoError(t, theAbility.ToFlush())
		})

		t.Run("fails when the URL is empty", func(t *testing.T) {
			require.ErrorIs(t, abilityWithoutConnection().ToConnectTo(""), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when the URL is unreachable", func(t *testing.T) {
			require.ErrorIs(
				t,
				abilityWithoutConnection().ToConnectTo("nats://invalid:4222"),
				ability.ErrFailedToConnect,
			)
		})

		t.Run("fails when already connected", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.ErrorIs(t, theAbility.ToConnectTo(nc.ConnectedUrl()), ability.ErrAlreadyConnected)
		})
	})

	t.Run("reports what happened to the connection", func(t *testing.T) {
		t.Run("nothing yet", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			connection, err := theAbility.ToRetrieveTheConnection()
			require.NoError(t, err)
			assert.False(t, connection.WasClosed())
			assert.False(t, connection.WasDisconnected())
			assert.False(t, connection.WasReconnected())
		})

		t.Run("closed", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			nc.Close()
			assertEventually(t, func() bool {
				connection, err := theAbility.ToRetrieveTheConnection()
				return err == nil && connection.WasClosed()
			})
		})

		t.Run("disconnected and reconnected", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.NoError(t, nc.ForceReconnect())
			assertEventually(t, func() bool {
				connection, err := theAbility.ToRetrieveTheConnection()
				return err == nil && connection.WasDisconnected() && connection.WasReconnected()
			})
			connection, err := theAbility.ToRetrieveTheConnection()
			require.NoError(t, err)
			assert.False(t, connection.WasClosed())
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			_, err := abilityWithoutConnection().ToRetrieveTheConnection()
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})
	})

	t.Run("publishes messages", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			sub, err := nc.SubscribeSync(subject)
			require.NoError(t, err)

			require.NoError(t, theAbility.ToPublish(ability.Message().To(subject).Containing(message)))

			msg, err := sub.NextMsg(testTimeout)
			require.NoError(t, err)
			assert.Equal(t, message, msg.Data)
			assert.Equal(t, subject, msg.Subject)
		})

		t.Run("fails when the message is nil", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.ErrorIs(t, theAbility.ToPublish(nil), ability.ErrMessageIsNil)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.ErrorIs(t, theAbility.ToPublish(ability.Message().Containing(message)), nats.ErrBadSubject)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			err := abilityWithoutConnection().ToPublish(ability.Message().To(subject))
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})
	})

	t.Run("flushes the connection", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.NoError(t, theAbility.ToPublish(ability.Message().To(subject).Containing(message)))
			require.NoError(t, theAbility.ToFlush())
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			require.ErrorIs(t, abilityWithoutConnection().ToFlush(), ability.ErrNATSConnectionNotSet)
		})
	})

	t.Run("sends requests", func(t *testing.T) {
		t.Run("successfully, and stores the reply", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
				assert.NoError(t, msg.Respond([]byte("ok")))
			})
			require.NoError(t, err)
			defer func() { assert.NoError(t, sub.Unsubscribe()) }()

			request := ability.Message().To(subject).Containing(message)
			require.NoError(t, theAbility.ToSendTheRequest(ctx, request, requestTimeout))

			replies, err := theAbility.ToRetrieveMessagesFrom(subject)
			require.NoError(t, err)
			require.Len(t, replies, 1)
			assert.Equal(t, []byte("ok"), replies[0].Data)
		})

		t.Run("fails when the message is nil", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.ErrorIs(t, theAbility.ToSendTheRequest(ctx, nil, requestTimeout), ability.ErrMessageIsNil)
		})

		t.Run("fails when the subject is empty", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			err := theAbility.ToSendTheRequest(ctx, ability.Message().Containing(message), requestTimeout)
			require.ErrorIs(t, err, nats.ErrBadSubject)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			err := abilityWithoutConnection().ToSendTheRequest(ctx, ability.Message().To(subject), requestTimeout)
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when no responder is set", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			err := theAbility.ToSendTheRequest(ctx, ability.Message().To(subject), requestTimeout)
			require.ErrorIs(t, err, nats.ErrNoResponders)
		})

		t.Run("fails when the reply takes too long", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			sub, err := nc.Subscribe(subject, func(_ *nats.Msg) {})
			require.NoError(t, err)
			defer func() { assert.NoError(t, sub.Unsubscribe()) }()

			err = theAbility.ToSendTheRequest(ctx, ability.Message().To(subject), 20*time.Millisecond)
			require.ErrorIs(t, err, nats.ErrTimeout)
			replies, err := theAbility.ToRetrieveMessagesFrom(subject)
			require.NoError(t, err)
			assert.Empty(t, replies)
		})

		t.Run("stops when the context ends", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			sub, err := nc.Subscribe(subject, func(_ *nats.Msg) {})
			require.NoError(t, err)
			defer func() { assert.NoError(t, sub.Unsubscribe()) }()

			cancelled, cancel := context.WithCancel(ctx)
			cancel()
			err = theAbility.ToSendTheRequest(cancelled, ability.Message().To(subject), requestTimeout)
			require.ErrorIs(t, err, context.Canceled)
		})
	})

	t.Run("subscribes to a subject", func(t *testing.T) {
		testSubscriptions(ctx, t, message)

		t.Run("fails when the queue is empty", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			err := theAbility.ToSubscribe(ability.Subscription{Subject: subject, Queue: " "})
			require.ErrorIs(t, err, nats.ErrBadQueueName)
		})

		t.Run("stores the messages of an asynchronous subscription", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject, Asynchronous: true}))
			require.NoError(t, theAbility.ToPublish(ability.Message().To(subject).Containing(message)))

			assertEventually(t, func() bool {
				messages, err := theAbility.ToRetrieveMessagesFrom(subject)
				return err == nil && len(messages) == 1 && string(messages[0].Data) == string(message)
			})
		})
	})

	t.Run("unsubscribes from a subject", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))
			assert.Equal(t, 1, nc.NumSubscriptions())

			require.NoError(t, theAbility.ToUnsubscribeFrom(subject))
			assert.Equal(t, 0, nc.NumSubscriptions())
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			require.ErrorIs(t, abilityWithoutConnection().ToUnsubscribeFrom(subject), ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when there is no subscription for the subject", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.ErrorIs(t, theAbility.ToUnsubscribeFrom(subject), ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the connection was closed", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))
			nc.Close()
			require.Error(t, theAbility.ToUnsubscribeFrom(subject))
		})
	})

	t.Run("receives the next message from a subject", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))
			require.NoError(t, nc.Publish(subject, message))

			require.NoError(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject))

			messages, err := theAbility.ToRetrieveMessagesFrom(subject)
			require.NoError(t, err)
			require.Len(t, messages, 1)
			assert.Equal(t, message, messages[0].Data)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			err := abilityWithoutConnection().ToReceiveTheNextMessageFrom(ctx, subject)
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})

		t.Run("fails when there is no subscription for the subject", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.ErrorIs(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject), ability.ErrNoSubscriptionFound)
		})

		t.Run("fails when the subscription is asynchronous", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject, Asynchronous: true}))
			require.ErrorIs(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject), nats.ErrSyncSubRequired)
		})

		t.Run("stops when the context ends", func(t *testing.T) {
			theAbility, _ := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))
			expiring, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
			defer cancel()

			err := theAbility.ToReceiveTheNextMessageFrom(expiring, subject)
			require.ErrorIs(t, err, context.DeadlineExceeded)
		})

		t.Run("fails when the connection was closed", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))
			nc.Close()
			require.Error(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject))
		})
	})

	t.Run("retrieves the messages", func(t *testing.T) {
		t.Run("in the order they were received, as a copy", func(t *testing.T) {
			theAbility, nc := connectedAbility(t)
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: "first"}))
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: "second"}))
			require.NoError(t, nc.Publish("first", []byte("1")))
			require.NoError(t, nc.Publish("second", []byte("2")))
			require.NoError(t, nc.Publish("first", []byte("3")))
			require.NoError(t, theAbility.ToReceiveTheNextMessageFrom(ctx, "first"))
			require.NoError(t, theAbility.ToReceiveTheNextMessageFrom(ctx, "second"))
			require.NoError(t, theAbility.ToReceiveTheNextMessageFrom(ctx, "first"))

			messages, err := theAbility.ToRetrieveMessages()
			require.NoError(t, err)
			require.Len(t, messages, 3)
			assert.Equal(t, []byte("1"), messages[0].Data)
			assert.Equal(t, []byte("2"), messages[1].Data)
			assert.Equal(t, []byte("3"), messages[2].Data)

			messages[0] = nil
			again, err := theAbility.ToRetrieveMessages()
			require.NoError(t, err)
			assert.NotNil(t, again[0])

			fromFirst, err := theAbility.ToRetrieveMessagesFrom("first")
			require.NoError(t, err)
			require.Len(t, fromFirst, 2)
			assert.Equal(t, []byte("1"), fromFirst[0].Data)
			assert.Equal(t, []byte("3"), fromFirst[1].Data)

			fromUnknown, err := theAbility.ToRetrieveMessagesFrom("unknown")
			require.NoError(t, err)
			assert.Empty(t, fromUnknown)
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			_, err := abilityWithoutConnection().ToRetrieveMessages()
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
			_, err = abilityWithoutConnection().ToRetrieveMessagesFrom(subject)
			require.ErrorIs(t, err, ability.ErrNATSConnectionNotSet)
		})
	})

	t.Run("forgets its resources", func(t *testing.T) {
		t.Run("closes the connection and drops the subscriptions and the messages", func(t *testing.T) {
			nc := fixture.CreateEmbeddedNATSServer(t)
			communicateWithNATS := ability.CommunicateWithNATS().Using(nc)
			theAbility := ability.TheAbilityOf(screenplay.ActorNamed("Adam").WhoCan(communicateWithNATS))
			require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject}))
			require.NoError(t, nc.Publish(subject, message))
			require.NoError(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject))

			require.NoError(t, communicateWithNATS.Forget())

			assert.True(t, nc.IsClosed())
			assert.Equal(t, 0, nc.NumSubscriptions())
			require.ErrorIs(t, theAbility.ToFlush(), ability.ErrNATSConnectionNotSet)
		})

		t.Run("does nothing when the connection is not set", func(t *testing.T) {
			require.NoError(t, ability.CommunicateWithNATS().Forget())
		})

		t.Run("fails on a nil receiver", func(t *testing.T) {
			var communicateWithNATS *ability.CommunicateWithNATSAbility
			require.ErrorIs(t, communicateWithNATS.Forget(), screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is described as a string", func(t *testing.T) {
		assert.Equal(t, "communicate with NATS", ability.CommunicateWithNATS().String())
	})

	t.Run("is safe for concurrent use", func(t *testing.T) {
		testConcurrentUse(t, message)
	})

	t.Run("fails every operation when the actor does not have the ability", func(t *testing.T) {
		theAbility := ability.TheAbilityOf(screenplay.ActorNamed("Adam"))
		require.ErrorIs(t, theAbility.ToConnectTo("nats://localhost:4222"), screenplay.ErrActorMissingAbility)
		require.ErrorIs(t, theAbility.ToPublish(ability.Message().To(subject)), screenplay.ErrActorMissingAbility)
		require.ErrorIs(
			t,
			theAbility.ToSendTheRequest(ctx, ability.Message().To(subject), requestTimeout),
			screenplay.ErrActorMissingAbility,
		)
		require.ErrorIs(t, theAbility.ToFlush(), screenplay.ErrActorMissingAbility)
		require.ErrorIs(
			t,
			theAbility.ToSubscribe(ability.Subscription{Subject: subject}),
			screenplay.ErrActorMissingAbility,
		)
		require.ErrorIs(t, theAbility.ToUnsubscribeFrom(subject), screenplay.ErrActorMissingAbility)
		require.ErrorIs(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject), screenplay.ErrActorMissingAbility)
		_, err := theAbility.ToRetrieveMessages()
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = theAbility.ToRetrieveMessagesFrom(subject)
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = theAbility.ToRetrieveTheConnection()
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("fails every operation when the actor is nil", func(t *testing.T) {
		require.ErrorIs(t, ability.TheAbilityOf(nil).ToFlush(), screenplay.ErrActorMissingAbility)
	})
}

// testSubscriptions checks every kind of subscription: synchronous or
// asynchronous, with or without a queue group.
func testSubscriptions(ctx context.Context, t *testing.T, message []byte) {
	t.Helper()

	cases := []struct {
		name         string
		subscription ability.Subscription
	}{
		{"synchronously", ability.Subscription{Subject: subject}},
		{"asynchronously", ability.Subscription{Subject: subject, Asynchronous: true}},
		{"synchronously with a queue", ability.Subscription{Subject: subject, Queue: queue}},
		{"asynchronously with a queue", ability.Subscription{Subject: subject, Queue: queue, Asynchronous: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("successfully", func(t *testing.T) {
				theAbility, nc := connectedAbility(t)
				require.NoError(t, theAbility.ToSubscribe(tc.subscription))
				assert.Equal(t, 1, nc.NumSubscriptions())
			})

			t.Run("unsubscribes after the wanted number of messages", func(t *testing.T) {
				theAbility, nc := connectedAbility(t)
				subscription := tc.subscription
				subscription.MaxMessages = 2
				require.NoError(t, theAbility.ToSubscribe(subscription))

				for range 3 {
					require.NoError(t, theAbility.ToPublish(ability.Message().To(subject).Containing(message)))
				}
				if !tc.subscription.Asynchronous {
					require.NoError(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject))
					require.NoError(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject))
					require.ErrorIs(t, theAbility.ToReceiveTheNextMessageFrom(ctx, subject), nats.ErrMaxMessages)
				}
				assertEventually(t, func() bool { return nc.NumSubscriptions() == 0 })
				messages, err := theAbility.ToRetrieveMessagesFrom(subject)
				require.NoError(t, err)
				assert.Len(t, messages, 2)
			})

			t.Run("fails when the subject is empty", func(t *testing.T) {
				theAbility, _ := connectedAbility(t)
				subscription := tc.subscription
				subscription.Subject = ""
				require.ErrorIs(t, theAbility.ToSubscribe(subscription), nats.ErrBadSubject)
			})

			t.Run("fails when the connection is not set", func(t *testing.T) {
				require.ErrorIs(
					t,
					abilityWithoutConnection().ToSubscribe(tc.subscription),
					ability.ErrNATSConnectionNotSet,
				)
			})

			t.Run("fails when the subscription already exists", func(t *testing.T) {
				theAbility, _ := connectedAbility(t)
				require.NoError(t, theAbility.ToSubscribe(tc.subscription))
				require.ErrorIs(t, theAbility.ToSubscribe(tc.subscription), ability.ErrSubscriptionAlreadyExists)
			})

			t.Run("fails when the connection is closed", func(t *testing.T) {
				theAbility, nc := connectedAbility(t)
				nc.Close()
				require.Error(t, theAbility.ToSubscribe(tc.subscription))
			})
		})
	}
}

// testConcurrentUse publishes, receives and retrieves messages at the same
// time; run with -race.
func testConcurrentUse(t *testing.T, message []byte) {
	t.Helper()

	theAbility, _ := connectedAbility(t)
	require.NoError(t, theAbility.ToSubscribe(ability.Subscription{Subject: subject, Asynchronous: true}))

	const workers = 10
	errs := make(chan error, 4*workers)
	var wg sync.WaitGroup
	for i := range workers {
		wg.Go(func() {
			errs <- theAbility.ToPublish(ability.Message().To(subject).Containing(message))
			_, err := theAbility.ToRetrieveMessages()
			errs <- err
			_, err = theAbility.ToRetrieveTheConnection()
			errs <- err
			if i%2 == 0 {
				errs <- theAbility.ToFlush()
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	assertEventually(t, func() bool {
		messages, err := theAbility.ToRetrieveMessagesFrom(subject)
		return err == nil && len(messages) == workers
	})
}

// connectedAbility returns the operations of an actor whose ability holds a
// connection to an embedded NATS server, and that connection.
func connectedAbility(t *testing.T) (*ability.Operations, *nats.Conn) {
	t.Helper()

	nc := fixture.CreateEmbeddedNATSServer(t)
	adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
	t.Cleanup(func() { assert.NoError(t, adam.Exit()) })

	return ability.TheAbilityOf(adam), nc
}

// abilityWithoutConnection returns the operations of an actor whose ability
// has neither a connection nor the URL of a server.
func abilityWithoutConnection() *ability.Operations {
	return ability.TheAbilityOf(screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS()))
}

func assertEventually(t *testing.T, condition func() bool) {
	t.Helper()
	assert.Eventually(t, condition, testTimeout, tickInterval)
}
