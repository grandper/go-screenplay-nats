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

	// ErrSubjectIsEmpty is returned when the subject is empty.
	ErrSubjectIsEmpty = errors.New("subject is empty")

	// ErrQueueIsEmpty is returned when the queue is empty.
	ErrQueueIsEmpty = errors.New("queue is empty")

	// ErrMessageIsNil is returned when the message is nil.
	ErrMessageIsNil = errors.New("message is nil")

	// ErrSubscriptionAlreadyExists is returned when a subscription for the subject already exists.
	ErrSubscriptionAlreadyExists = errors.New("subscription for the subject already exists")

	// ErrNoSubscriptionFound is returned when no subscription is found for the subject.
	ErrNoSubscriptionFound = errors.New("no subscription found for the subject")

	// ErrNoMessageForSubject is returned when no message is found for the subject.
	ErrNoMessageForSubject = errors.New("no message found for the subject")
)

// CommunicateWithNATSAbility represents the ability to communicate with a NATS server.
// It is safe for concurrent use.
type CommunicateWithNATSAbility struct {
	// connecting serializes the lazy connection so that the state mutex is
	// never held during a network call.
	connecting sync.Mutex
	url        string

	mutex            sync.RWMutex
	nc               *nats.Conn
	subs             map[string]*nats.Subscription
	receivedMessages map[string][]*nats.Msg

	connWasClosed bool
	connClosedErr error

	asyncErrOccurred bool
	asyncErr         error

	connDisconnected    bool
	connDisconnectedErr error

	connReconnected bool
}

// CommunicateWithNATS creates a new instance of the CommunicateWithNATS ability.
// This ability allows the user to communicate with a NATS server. The
// connection is given with Using, or opened by the ability itself from the URL
// given with UsingServerAt.
func CommunicateWithNATS() *CommunicateWithNATSAbility {
	return &CommunicateWithNATSAbility{
		subs:             make(map[string]*nats.Subscription),
		receivedMessages: make(map[string][]*nats.Msg),
	}
}

// Using provides the NATS connection to the ability. A connection already held
// by the ability is kept, and a nil connection is ignored.
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

// UsingServerAt gives the ability the URL of the NATS server to connect to. The
// connection is opened the first time the ability is used, so that a failure to
// connect is reported by the operation that needed it. A connection already
// held by the ability is kept.
func (c *CommunicateWithNATSAbility) UsingServerAt(url string) *CommunicateWithNATSAbility {
	if c == nil {
		return c
	}

	c.connecting.Lock()
	defer c.connecting.Unlock()

	c.url = url

	return c
}

// ConnectTo connects the ability to the NATS server at the given URL. It fails
// with ErrAlreadyConnected when the ability already holds a connection, and
// with an error wrapping ErrFailedToConnect when the server cannot be reached.
func (c *CommunicateWithNATSAbility) ConnectTo(url string) error {
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

// use stores the connection and records what happens to it. The state mutex
// must be held.
func (c *CommunicateWithNATSAbility) use(nc *nats.Conn) {
	c.nc = nc
	c.setupHandlers()
}

// connect opens a connection to the server at the given URL and stores it. The
// connecting mutex must be held.
func (c *CommunicateWithNATSAbility) connect(url string) error {
	if url == "" {
		return ErrNATSConnectionNotSet
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return fmt.Errorf("%w at %s: %w", ErrFailedToConnect, url, err)
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

// connection returns the connection of the ability, opening it from the URL
// given with UsingServerAt when needed. It fails with
// screenplay.ErrActorMissingAbility on a nil ability and with
// ErrNATSConnectionNotSet when the ability has neither a connection nor a URL.
func (c *CommunicateWithNATSAbility) connection() (*nats.Conn, error) {
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

func (c *CommunicateWithNATSAbility) setupHandlers() {
	c.nc.SetClosedHandler(func(conn *nats.Conn) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.connWasClosed = true
		c.connClosedErr = conn.LastError()
	})

	c.nc.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.asyncErrOccurred = true
		c.asyncErr = err
	})

	c.nc.SetDisconnectErrHandler(func(_ *nats.Conn, err error) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.connDisconnected = true
		c.connDisconnectedErr = err
	})

	c.nc.SetReconnectHandler(func(_ *nats.Conn) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.connReconnected = true
	})
}

// Publish publishes a message.
func (c *CommunicateWithNATSAbility) Publish(message *NATSMessage) error {
	nc, err := c.connection()
	if err != nil {
		return err
	}

	if message == nil {
		return ErrMessageIsNil
	}

	if message.Subject() == "" {
		return ErrSubjectIsEmpty
	}

	return nc.PublishMsg(message.toNATSMessage())
}

// Request sends a request.
func (c *CommunicateWithNATSAbility) Request(message *NATSMessage, requestTimeout time.Duration) (*nats.Msg, error) {
	nc, err := c.connection()
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, ErrMessageIsNil
	}

	if message.Subject() == "" {
		return nil, ErrSubjectIsEmpty
	}

	response, err := nc.RequestMsg(message.toNATSMessage(), requestTimeout)
	if err != nil {
		return nil, err
	}
	c.store(message.Subject(), response)
	return response, nil
}

// Respond sends a response to a message.
func (c *CommunicateWithNATSAbility) Respond(m *nats.Msg, response string) error {
	if _, err := c.connection(); err != nil {
		return err
	}
	if m == nil {
		return ErrMessageIsNil
	}
	return m.Respond([]byte(response))
}

// Flush forces a flush of the buffered publish operations on the NATS connection.
func (c *CommunicateWithNATSAbility) Flush() error {
	nc, err := c.connection()
	if err != nil {
		return err
	}
	return nc.Flush()
}

// SubscribeAsyncTo subscribes asynchronously to a subject.
func (c *CommunicateWithNATSAbility) SubscribeAsyncTo(subject string, maxWantedMessage int) error {
	return c.subscribe(subject, "", maxWantedMessage, func(nc *nats.Conn) (*nats.Subscription, error) {
		return nc.Subscribe(subject, c.handlerFor(subject))
	})
}

// SubscribeTo subscribes synchronously to a subject.
func (c *CommunicateWithNATSAbility) SubscribeTo(subject string, maxWantedMessage int) error {
	return c.subscribe(subject, "", maxWantedMessage, func(nc *nats.Conn) (*nats.Subscription, error) {
		return nc.SubscribeSync(subject)
	})
}

// SubscribeToQueueAsync subscribes asynchronously to a subject with a queue group.
// All subscriptions with the same queue name form a queue group, and each
// message is delivered to only one subscriber of the group.
func (c *CommunicateWithNATSAbility) SubscribeToQueueAsync(subject, queue string, maxWantedMessage int) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}
	if queue == "" {
		return ErrQueueIsEmpty
	}
	return c.subscribe(subject, queue, maxWantedMessage, func(nc *nats.Conn) (*nats.Subscription, error) {
		return nc.QueueSubscribe(subject, queue, c.handlerFor(subject))
	})
}

// SubscribeToQueue subscribes synchronously to a subject with a queue group.
func (c *CommunicateWithNATSAbility) SubscribeToQueue(subject, queue string, maxWantedMessage int) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}
	if queue == "" {
		return ErrQueueIsEmpty
	}
	return c.subscribe(subject, queue, maxWantedMessage, func(nc *nats.Conn) (*nats.Subscription, error) {
		return nc.QueueSubscribeSync(subject, queue)
	})
}

// subscribe registers the subscription created by subscribeWith under the
// subject, and bounds it to maxWantedMessage messages when that number is
// positive. The queue is only used to describe a failure.
func (c *CommunicateWithNATSAbility) subscribe(
	subject, queue string,
	maxWantedMessage int,
	subscribeWith func(nc *nats.Conn) (*nats.Subscription, error),
) error {
	if c == nil {
		return screenplay.ErrActorMissingAbility
	}

	if subject == "" {
		return ErrSubjectIsEmpty
	}

	nc, err := c.connection()
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.subs[subject]; exists {
		return ErrSubscriptionAlreadyExists
	}

	sub, err := subscribeWith(nc)
	if err != nil {
		if queue != "" {
			return fmt.Errorf("failed to subscribe to subject %s with queue %s: %w", subject, queue, err)
		}
		return err
	}

	c.subs[subject] = sub

	if maxWantedMessage <= 0 {
		return nil
	}
	return sub.AutoUnsubscribe(maxWantedMessage)
}

// handlerFor records every message received on the subject.
func (c *CommunicateWithNATSAbility) handlerFor(subject string) nats.MsgHandler {
	return func(msg *nats.Msg) {
		c.store(subject, msg)
	}
}

// store records a message received on the subject.
func (c *CommunicateWithNATSAbility) store(subject string, msg *nats.Msg) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.receivedMessages[subject] = append(c.receivedMessages[subject], msg)
}

// UnsubscribeFrom unsubscribes from a subject.
func (c *CommunicateWithNATSAbility) UnsubscribeFrom(subject string) error {
	if c == nil {
		return screenplay.ErrActorMissingAbility
	}

	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if _, err := c.connection(); err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	sub, ok := c.subs[subject]
	if !ok {
		return fmt.Errorf("%w %s", ErrNoSubscriptionFound, subject)
	}

	if err := sub.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe from subject %s: %w", subject, err)
	}

	delete(c.subs, subject)
	return nil
}

// subscription returns the subscription registered under the subject.
func (c *CommunicateWithNATSAbility) subscription(subject string) (*nats.Subscription, error) {
	if c == nil {
		return nil, screenplay.ErrActorMissingAbility
	}

	if subject == "" {
		return nil, ErrSubjectIsEmpty
	}

	if _, err := c.connection(); err != nil {
		return nil, err
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	sub, ok := c.subs[subject]
	if !ok {
		return nil, fmt.Errorf("%w %s", ErrNoSubscriptionFound, subject)
	}

	return sub, nil
}

// ReceiveNextMessageFor receives the next message for the specified subject.
func (c *CommunicateWithNATSAbility) ReceiveNextMessageFor(ctx context.Context, subject string) error {
	sub, err := c.subscription(subject)
	if err != nil {
		return err
	}

	msg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next msg from NATS for subject %s: %w", subject, err)
	}

	c.store(subject, msg)
	return nil
}

// ReceiveAllMessagesFor receives all messages for the specified subject until the context is done.
func (c *CommunicateWithNATSAbility) ReceiveAllMessagesFor(
	ctx context.Context,
	subject string,
	timeout time.Duration,
) error {
	sub, err := c.subscription(subject)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		msg, errNext := sub.NextMsgWithContext(ctx)
		if errors.Is(errNext, context.DeadlineExceeded) {
			break
		}
		if errNext != nil {
			return fmt.Errorf("failed to get next msg from NATS for subject %s: %w", subject, errNext)
		}
		c.store(subject, msg)
	}
	return nil
}

// ReceiveAllMessages receives all the messages from the subscribers.
func (c *CommunicateWithNATSAbility) ReceiveAllMessages(ctx context.Context, timeout time.Duration) error {
	if _, err := c.connection(); err != nil {
		return err
	}

	for _, subject := range c.subjects() {
		if err := c.ReceiveAllMessagesFor(ctx, subject, timeout); err != nil {
			return err
		}
	}
	return nil
}

// subjects returns the subjects the ability is subscribed to.
func (c *CommunicateWithNATSAbility) subjects() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	subjects := make([]string, 0, len(c.subs))
	for subject := range c.subs {
		subjects = append(subjects, subject)
	}
	return subjects
}

// LastMessageForSubject retrieves the last message received for the specified subject.
func (c *CommunicateWithNATSAbility) LastMessageForSubject(subject string) (*nats.Msg, error) {
	messages, err := c.AllMessagesForSubject(subject)
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("%w %s", ErrNoMessageForSubject, subject)
	}

	return messages[len(messages)-1], nil
}

// AllMessagesForSubject retrieves all the messages for the specified subject.
func (c *CommunicateWithNATSAbility) AllMessagesForSubject(subject string) ([]*nats.Msg, error) {
	if c == nil {
		return nil, screenplay.ErrActorMissingAbility
	}

	if subject == "" {
		return nil, ErrSubjectIsEmpty
	}

	if _, err := c.connection(); err != nil {
		return nil, err
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	messages, ok := c.receivedMessages[subject]
	if !ok {
		return nil, nil
	}
	return append([]*nats.Msg(nil), messages...), nil
}

// AllMessages retrieves all the messages.
func (c *CommunicateWithNATSAbility) AllMessages() ([]*nats.Msg, error) {
	if _, err := c.connection(); err != nil {
		return nil, err
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	numMessages := c.allMessageCount()
	allMessages := make([]*nats.Msg, 0, numMessages)
	for _, messages := range c.receivedMessages {
		allMessages = append(allMessages, messages...)
	}
	return allMessages, nil
}

func (c *CommunicateWithNATSAbility) allMessageCount() int {
	count := 0
	for _, messages := range c.receivedMessages {
		count += len(messages)
	}
	return count
}

// ConnectionWasClosed indicates whether the NATS connection was closed.
func (c *CommunicateWithNATSAbility) ConnectionWasClosed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.connWasClosed
}

// ConnectionAsyncErrorOccurred indicates whether an error occurred on the NATS connection.
func (c *CommunicateWithNATSAbility) ConnectionAsyncErrorOccurred() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.asyncErrOccurred
}

// ConnectionWasDisconnected indicates whether the NATS connection was disconnected.
func (c *CommunicateWithNATSAbility) ConnectionWasDisconnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.connDisconnected
}

// ConnectionWasReconnected indicates whether the NATS connection was reconnected.
func (c *CommunicateWithNATSAbility) ConnectionWasReconnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.connReconnected
}

// Forget cleans up the ability's resources: every subscription is dropped and
// the connection is closed.
func (c *CommunicateWithNATSAbility) Forget() error {
	if c == nil {
		return nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.nc == nil {
		return nil
	}
	for _, sub := range c.subs {
		_ = sub.Unsubscribe()
	}
	c.subs = map[string]*nats.Subscription{}
	c.nc.Close()
	return nil
}

// CommunicateWithNATSAbility implements the screenplay.Ability interface.
var _ screenplay.Ability = (*CommunicateWithNATSAbility)(nil)
