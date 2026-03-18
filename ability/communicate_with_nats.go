package ability

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	// ErrNATSConnectionNotSet is returned when the NATS connection is not set.
	ErrNATSConnectionNotSet = errors.New("nats connection not set")

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
type CommunicateWithNATSAbility struct {
	nc               *nats.Conn
	subs             map[string]*nats.Subscription
	receivedMessages map[string][]*nats.Msg

	mutex sync.RWMutex

	connWasClosed bool
	connClosedErr error

	asyncErrOccurred bool
	asyncErr         error

	connDisconnected    bool
	connDisconnectedErr error

	connReconnected bool
}

// CommunicateWithNATS creates a new instance of the CommunicateWithNATS ability.
// This ability allows the user to communicate with a NATS server.
func CommunicateWithNATS() *CommunicateWithNATSAbility {
	return nil
}

// Using provides the NATS connection to the ability.
func (c *CommunicateWithNATSAbility) Using(nc *nats.Conn) *CommunicateWithNATSAbility {
	if nc == nil {
		return nil
	}
	communicateWithNATS := &CommunicateWithNATSAbility{
		nc:               nc,
		subs:             make(map[string]*nats.Subscription),
		receivedMessages: make(map[string][]*nats.Msg),
	}
	communicateWithNATS.nc = nc
	communicateWithNATS.setupHandlers()
	return communicateWithNATS
}

// UsingServerAt establishes a connection to the NATS server at the specified URL.
func (c *CommunicateWithNATSAbility) UsingServerAt(url string) *CommunicateWithNATSAbility {
	if url == "" {
		return nil
	}
	nc, err := nats.Connect(url)
	if err != nil {
		return nil
	}
	return c.Using(nc)
}

func (c *CommunicateWithNATSAbility) setupHandlers() {
	c.nc.SetClosedHandler(func(conn *nats.Conn) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.connWasClosed = true
		c.connClosedErr = conn.LastError()
	})

	c.nc.SetErrorHandler(func(_ *nats.Conn, subscription *nats.Subscription, err error) {
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

	c.nc.SetReconnectHandler(func(conn *nats.Conn) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.connReconnected = true
	})
}

// Publish publishes a message.
func (c *CommunicateWithNATSAbility) Publish(message *NATSMessage) error {
	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	if message == nil {
		return ErrMessageIsNil
	}

	if message.Subject() == "" {
		return ErrSubjectIsEmpty
	}

	return c.nc.PublishMsg(message.toNATSMessage())
}

// Request sends a request.
func (c *CommunicateWithNATSAbility) Request(message *NATSMessage, requestTimeout time.Duration) (*nats.Msg, error) {
	if c == nil || c.nc == nil {
		return nil, ErrNATSConnectionNotSet
	}

	if message == nil {
		return nil, ErrMessageIsNil
	}

	if message.Subject() == "" {
		return nil, ErrSubjectIsEmpty
	}

	return c.nc.RequestMsg(message.toNATSMessage(), requestTimeout)
}

// Respond sends a response to a message.
func (c *CommunicateWithNATSAbility) Respond(m *nats.Msg, response string) error {
	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}
	if m == nil {
		return ErrMessageIsNil
	}
	return m.Respond([]byte(response))
}

// SubscribeAsyncTo subscribes asynchronously to a subject.
func (c *CommunicateWithNATSAbility) SubscribeAsyncTo(subject string, maxWantedMessage int) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	if _, exists := c.subs[subject]; exists {
		return ErrSubscriptionAlreadyExists
	}

	sub, err := c.nc.Subscribe(subject, func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		return err
	}

	c.subs[subject] = sub

	if maxWantedMessage <= 0 {
		return nil
	}
	return sub.AutoUnsubscribe(maxWantedMessage)
}

// SubscribeTo subscribes synchronously to a subject.
func (c *CommunicateWithNATSAbility) SubscribeTo(subject string, maxWantedMessage int) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	if _, exists := c.subs[subject]; exists {
		return ErrSubscriptionAlreadyExists
	}

	sub, err := c.nc.SubscribeSync(subject)
	if err != nil {
		return err
	}

	c.subs[subject] = sub

	if maxWantedMessage <= 0 {
		return nil
	}
	return sub.AutoUnsubscribe(maxWantedMessage)
}

// SubscribeToQueueAsync subscribes asynchronously to a subject with a queue group.
func (c *CommunicateWithNATSAbility) SubscribeToQueueAsync(subject, queue string, maxWantedMessage int) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if queue == "" {
		return ErrQueueIsEmpty
	}

	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	if _, exists := c.subs[subject]; exists {
		return ErrSubscriptionAlreadyExists
	}

	// Queue Subscriber
	// All subscriptions with the same queue name will form a queue group.
	// Each message will be delivered to only one subscriber per queue group,
	// using queuing semantics. You can have as many queue groups as you wish.
	// Normal subscribers will continue to work as expected.
	sub, err := c.nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		c.mutex.Lock()
		c.receivedMessages[subject] = append(c.receivedMessages[subject], msg)
		c.mutex.Unlock()
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s with queue %s: %w", subject, queue, err)
	}

	c.subs[subject] = sub

	if maxWantedMessage <= 0 {
		return nil
	}
	return sub.AutoUnsubscribe(maxWantedMessage)
}

// SubscribeToQueue subscribes synchronously to a subject with a queue group.
func (c *CommunicateWithNATSAbility) SubscribeToQueue(subject, queue string, maxWantedMessage int) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if queue == "" {
		return ErrQueueIsEmpty
	}

	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	if _, exists := c.subs[subject]; exists {
		return ErrSubscriptionAlreadyExists
	}

	sub, err := c.nc.QueueSubscribeSync(subject, queue)
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s with queue %s: %w", subject, queue, err)
	}

	c.subs[subject] = sub

	if maxWantedMessage <= 0 {
		return nil
	}
	return sub.AutoUnsubscribe(maxWantedMessage)
}

// UnsubscribeFrom unsubscribes from a subject.
func (c *CommunicateWithNATSAbility) UnsubscribeFrom(subject string) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

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

// ReceiveNextMessageFor receives the next message for the specified subject.
func (c *CommunicateWithNATSAbility) ReceiveNextMessageFor(ctx context.Context, subject string) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	sub, ok := c.subs[subject]
	if !ok {
		return fmt.Errorf("%w %s", ErrNoSubscriptionFound, subject)
	}

	msg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next msg from NATS for subject %s: %w", subject, err)
	}

	c.mutex.Lock()
	c.receivedMessages[subject] = append(c.receivedMessages[subject], msg)
	c.mutex.Unlock()
	return nil
}

// ReceiveAllMessagesFor receives all messages for the specified subject until the context is done.
func (c *CommunicateWithNATSAbility) ReceiveAllMessagesFor(ctx context.Context, subject string, timeout time.Duration) error {
	if subject == "" {
		return ErrSubjectIsEmpty
	}

	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	sub, ok := c.subs[subject]
	if !ok {
		return fmt.Errorf("%w %s", ErrNoSubscriptionFound, subject)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		msg, err := sub.NextMsgWithContext(ctx)
		if errors.Is(err, context.DeadlineExceeded) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to get next msg from NATS for subject %s: %w", subject, err)
		}
		c.mutex.Lock()
		c.receivedMessages[subject] = append(c.receivedMessages[subject], msg)
		c.mutex.Unlock()
	}
	return nil
}

// ReceiveAllMessages receives all the messages from the subscribers.
func (c *CommunicateWithNATSAbility) ReceiveAllMessages(ctx context.Context, timeout time.Duration) error {
	if c == nil || c.nc == nil {
		return ErrNATSConnectionNotSet
	}

	for subject := range c.subs {
		if err := c.ReceiveAllMessagesFor(ctx, subject, timeout); err != nil {
			return err
		}
	}
	return nil
}

// LastMessageForSubject retrieves the last message received for the specified subject.
func (c *CommunicateWithNATSAbility) LastMessageForSubject(subject string) (*nats.Msg, error) {
	if subject == "" {
		return nil, ErrSubjectIsEmpty
	}

	if c == nil || c.nc == nil {
		return nil, ErrNATSConnectionNotSet
	}

	c.mutex.RLock()
	messages := c.receivedMessages[subject]
	numMessages := len(messages)
	c.mutex.RUnlock()
	if numMessages == 0 {
		return nil, fmt.Errorf("%w %s", ErrNoMessageForSubject, subject)
	}

	return messages[numMessages-1], nil
}

// AllMessagesForSubject retrieves all the messages for the specified subject.
func (c *CommunicateWithNATSAbility) AllMessagesForSubject(subject string) ([]*nats.Msg, error) {
	if subject == "" {
		return nil, ErrSubjectIsEmpty
	}

	if c == nil || c.nc == nil {
		return nil, ErrNATSConnectionNotSet
	}

	c.mutex.RLock()
	messages, ok := c.receivedMessages[subject]
	c.mutex.RUnlock()
	if ok {
		return messages, nil
	}
	return nil, nil
}

// AllMessages retrieves all the messages.
func (c *CommunicateWithNATSAbility) AllMessages() ([]*nats.Msg, error) {
	if c == nil || c.nc == nil {
		return nil, ErrNATSConnectionNotSet
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

// Forget clean up the ability's resources.
func (c *CommunicateWithNATSAbility) Forget() error {
	return nil
}
