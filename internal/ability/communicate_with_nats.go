package ability

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/nats-io/nats.go"
)

var (
	// ErrNATSConnectionNotSet is returned when the ability has neither a NATS
	// connection nor the URL of a server to connect to.
	ErrNATSConnectionNotSet = errors.New("nats connection not set")

	// ErrFailedToConnect is returned when the connection to the NATS server fails.
	ErrFailedToConnect = errors.New("failed to connect to the NATS server")

	// ErrAlreadyConnected is returned when the ability is asked to connect while it
	// already holds a connection.
	ErrAlreadyConnected = errors.New("already connected to a NATS server")

	// ErrMessageIsNil is returned when the message to publish or to send is nil.
	ErrMessageIsNil = errors.New("message is nil")

	// ErrSubscriptionAlreadyExists is returned when a subscription for the subject already exists.
	ErrSubscriptionAlreadyExists = errors.New("subscription for the subject already exists")

	// ErrNoSubscriptionFound is returned when no subscription is found for the subject.
	ErrNoSubscriptionFound = errors.New("no subscription found for the subject")
)

// CommunicateWithNATS enables the actor to communicate with a NATS server. The
// connection is given with Using, or opened by the ability itself from the URL
// given with UsingServerAt, reading "communicate with NATS using ...":
//
//	CommunicateWithNATS().Using(nc)
//	CommunicateWithNATS().UsingServerAt("nats://localhost:4222")
func CommunicateWithNATS() *CommunicateWithNATSAbility {
	return &CommunicateWithNATSAbility{
		subscriptions: map[string]*nats.Subscription{},
	}
}

// CommunicateWithNATSAbility is the ability to communicate with a NATS server.
// It is safe for concurrent use: actions running in parallel (see
// proceed.ConcurrentlyTo) can publish, subscribe and receive without racing.
// Its method set holds only its settings and Forget; the operations are
// reached through TheAbilityOf, so a test cannot bypass the actor.
type CommunicateWithNATSAbility struct {
	// connecting serializes the lazy connection so that the state mutex is
	// never held during a network call.
	connecting sync.Mutex
	url        string

	mutex         sync.RWMutex
	nc            *nats.Conn
	subscriptions map[string]*nats.Subscription
	received      []receivedMessage

	wasClosed       bool
	wasDisconnected bool
	wasReconnected  bool
}

// receivedMessage is a message received through a subscription or as the reply
// to a request, with the subject it was received from.
type receivedMessage struct {
	subject string
	msg     *nats.Msg
}

// Using provides the NATS connection to the ability; it reads "using nc". A
// connection already held by the ability is kept, and a nil connection is
// ignored.
func (c *CommunicateWithNATSAbility) Using(nc *nats.Conn) *CommunicateWithNATSAbility {
	if c == nil || nc == nil {
		return c
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.nc == nil {
		c.use(nc)
	}

	return c
}

// UsingServerAt gives the ability the URL of the NATS server to connect to; it
// reads "using server at nats://localhost:4222". The connection is opened the
// first time the ability is used, so that a failure to connect is reported by
// the action that needed it. A connection already held by the ability is kept.
func (c *CommunicateWithNATSAbility) UsingServerAt(url string) *CommunicateWithNATSAbility {
	if c == nil {
		return c
	}

	c.connecting.Lock()
	defer c.connecting.Unlock()

	c.url = url

	return c
}

// Forget drops every subscription, forgets the received messages and closes
// the connection. On a nil receiver, the actor does not have the ability, it
// fails with screenplay.ErrActorMissingAbility.
func (c *CommunicateWithNATSAbility) Forget() error {
	if c == nil {
		return screenplay.ErrActorMissingAbility
	}

	c.connecting.Lock()
	defer c.connecting.Unlock()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.url = ""
	c.received = nil
	if c.nc == nil {
		return nil
	}
	for _, subscription := range c.subscriptions {
		_ = subscription.Unsubscribe()
	}
	c.subscriptions = map[string]*nats.Subscription{}
	c.nc.Close()
	c.nc = nil

	return nil
}

// String describes the ability.
func (c *CommunicateWithNATSAbility) String() string {
	return "communicate with NATS"
}

// use stores the connection and records what happens to it. The state mutex
// must be held.
func (c *CommunicateWithNATSAbility) use(nc *nats.Conn) {
	c.nc = nc

	nc.SetClosedHandler(func(_ *nats.Conn) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.wasClosed = true
	})
	nc.SetDisconnectErrHandler(func(_ *nats.Conn, _ error) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.wasDisconnected = true
	})
	nc.SetReconnectHandler(func(_ *nats.Conn) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.wasReconnected = true
	})
}

// CommunicateWithNATSAbility implements the screenplay.Ability interface.
var _ screenplay.Ability = (*CommunicateWithNATSAbility)(nil)

// TheAbilityOf returns the operations of the actor's ability to communicate
// with NATS; it is written in PerformAs or AnsweredBy and reads "the ability of
// the actor to publish":
//
//	err := ability.TheAbilityOf(theActor).ToPublish(message)
//
// It never fails. When the actor is nil or does not have the ability, every
// operation fails with screenplay.ErrActorMissingAbility.
func TheAbilityOf(actor *screenplay.Actor) *Operations {
	return &Operations{ability: screenplay.UseAnAbilityOf(actor).That(CommunicateWithNATS())}
}

// Operations holds the operations of the actor's ability to communicate with
// NATS. They are not methods of CommunicateWithNATSAbility, whose method set
// holds only its settings, so a test cannot publish a message without an action.
type Operations struct {
	ability *CommunicateWithNATSAbility
}

// Subscription describes a subscription to a subject: with a queue group when
// Queue is set, asynchronous when Asynchronous is set, and unsubscribed after
// MaxMessages messages when that number is positive.
type Subscription struct {
	Subject      string
	Queue        string
	Asynchronous bool
	MaxMessages  int
}

// ToConnectTo connects the ability to the NATS server at the given URL. It
// fails with ErrAlreadyConnected when the ability already holds a connection,
// and with an error wrapping ErrFailedToConnect when the server cannot be
// reached. When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToConnectTo(url string) error {
	c := o.ability
	if c == nil {
		return screenplay.ErrActorMissingAbility
	}

	c.connecting.Lock()
	defer c.connecting.Unlock()

	if c.isConnected() {
		return ErrAlreadyConnected
	}

	return c.connect(url)
}

// ToPublish publishes the message. When the actor does not have the ability,
// it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToPublish(message *NATSMessage) error {
	nc, err := o.connection()
	if err != nil {
		return err
	}
	if message == nil {
		return ErrMessageIsNil
	}

	return nc.PublishMsg(message.toNATSMessage())
}

// ToSendTheRequest sends the message as a request and stores the reply under
// the subject of the request. The request stops when the context ends or
// after the timeout, whichever comes first; a timeout is reported with an
// error wrapping nats.ErrTimeout. When the actor does not have the ability,
// it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToSendTheRequest(ctx context.Context, message *NATSMessage, timeout time.Duration) error {
	nc, err := o.connection()
	if err != nil {
		return err
	}
	if message == nil {
		return ErrMessageIsNil
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reply, err := nc.RequestMsgWithContext(ctx, message.toNATSMessage())
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf(
			"no reply to the request to %q within %s: %w: %w",
			message.Subject(),
			timeout,
			nats.ErrTimeout,
			err,
		)
	}
	if err != nil {
		return err
	}

	o.ability.store(message.Subject(), reply)

	return nil
}

// ToFlush forces a flush of the buffered publish operations on the NATS
// connection. When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToFlush() error {
	nc, err := o.connection()
	if err != nil {
		return err
	}

	return nc.Flush()
}

// ToSubscribe subscribes to the subject of the subscription. An asynchronous
// subscription stores every message it receives; a synchronous one hands its
// messages to ToReceiveTheNextMessageFrom. When the actor does not have the
// ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToSubscribe(subscription Subscription) error {
	nc, err := o.connection()
	if err != nil {
		return err
	}
	c := o.ability

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.subscriptions[subscription.Subject]; exists {
		return ErrSubscriptionAlreadyExists
	}

	sub, err := c.subscribe(nc, subscription)
	if err != nil {
		return err
	}
	c.subscriptions[subscription.Subject] = sub

	if subscription.MaxMessages <= 0 {
		return nil
	}

	return sub.AutoUnsubscribe(subscription.MaxMessages)
}

// subscribe creates the subscription on the connection. The state mutex must
// be held: the handler of an asynchronous subscription locks it only once a
// message arrives.
func (c *CommunicateWithNATSAbility) subscribe(nc *nats.Conn, subscription Subscription) (*nats.Subscription, error) {
	subject := subscription.Subject
	handler := func(msg *nats.Msg) {
		c.store(subject, msg)
	}

	switch {
	case subscription.Queue != "" && subscription.Asynchronous:
		return nc.QueueSubscribe(subject, subscription.Queue, handler)
	case subscription.Queue != "":
		return nc.QueueSubscribeSync(subject, subscription.Queue)
	case subscription.Asynchronous:
		return nc.Subscribe(subject, handler)
	default:
		return nc.SubscribeSync(subject)
	}
}

// ToUnsubscribeFrom unsubscribes from the subject. When the actor does not have
// the ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToUnsubscribeFrom(subject string) error {
	if _, err := o.connection(); err != nil {
		return err
	}
	c := o.ability

	c.mutex.Lock()
	defer c.mutex.Unlock()

	sub, ok := c.subscriptions[subject]
	if !ok {
		return fmt.Errorf("%w: %q", ErrNoSubscriptionFound, subject)
	}
	if err := sub.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe from %q: %w", subject, err)
	}
	delete(c.subscriptions, subject)

	return nil
}

// ToReceiveTheNextMessageFrom waits for the next message of the synchronous
// subscription to the subject, until the context ends, and stores it. When
// the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToReceiveTheNextMessageFrom(ctx context.Context, subject string) error {
	if _, err := o.connection(); err != nil {
		return err
	}
	c := o.ability

	c.mutex.RLock()
	sub, ok := c.subscriptions[subject]
	c.mutex.RUnlock()
	if !ok {
		return fmt.Errorf("%w: %q", ErrNoSubscriptionFound, subject)
	}

	msg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive the next message from %q: %w", subject, err)
	}
	c.store(subject, msg)

	return nil
}

// ToRetrieveMessages returns a copy of every message received, in the order
// they were received. When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToRetrieveMessages() ([]*nats.Msg, error) {
	return o.retrieveMessages(func(string) bool { return true })
}

// ToRetrieveMessagesFrom returns a copy of the messages received from the
// subject, in the order they were received. The subject is the one of the
// subscription, or of the request the message replied to. When the actor does
// not have the ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToRetrieveMessagesFrom(subject string) ([]*nats.Msg, error) {
	return o.retrieveMessages(func(received string) bool { return received == subject })
}

func (o *Operations) retrieveMessages(wanted func(subject string) bool) ([]*nats.Msg, error) {
	if _, err := o.connection(); err != nil {
		return nil, err
	}
	c := o.ability

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	messages := make([]*nats.Msg, 0, len(c.received))
	for _, received := range c.received {
		if wanted(received.subject) {
			messages = append(messages, received.msg)
		}
	}

	return messages, nil
}

// ToRetrieveTheConnection returns a snapshot of what happened to the
// connection. When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToRetrieveTheConnection() (*Connection, error) {
	if _, err := o.connection(); err != nil {
		return nil, err
	}
	c := o.ability

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return NewConnection(c.wasClosed, c.wasDisconnected, c.wasReconnected), nil
}

// connection returns the connection of the ability, opening it from the URL
// given with UsingServerAt when needed. It fails with
// screenplay.ErrActorMissingAbility when the actor does not have the ability,
// and with ErrNATSConnectionNotSet when the ability has neither a connection
// nor a URL.
func (o *Operations) connection() (*nats.Conn, error) {
	c := o.ability
	if c == nil {
		return nil, screenplay.ErrActorMissingAbility
	}

	c.connecting.Lock()
	defer c.connecting.Unlock()

	if !c.isConnected() {
		if err := c.connect(c.url); err != nil {
			return nil, err
		}
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.nc, nil
}

// connect opens a connection to the server at the given URL and stores it.
// The connecting mutex must be held.
func (c *CommunicateWithNATSAbility) connect(url string) error {
	if url == "" {
		return ErrNATSConnectionNotSet
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return fmt.Errorf("%w at %q: %w", ErrFailedToConnect, url, err)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.use(nc)

	return nil
}

func (c *CommunicateWithNATSAbility) isConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.nc != nil
}

// store records a message received from the subject.
func (c *CommunicateWithNATSAbility) store(subject string, msg *nats.Msg) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.received = append(c.received, receivedMessage{subject: subject, msg: msg})
}
