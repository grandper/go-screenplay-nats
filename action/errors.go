package action

import "errors"

var (
	// ErrURLIsEmpty is returned when connect.ToTheNATSServerAt is given an empty
	// URL. Nothing panics when the action is built; the mistake is reported when
	// the actor performs it.
	ErrURLIsEmpty = errors.New("the URL of the NATS server must be provided")

	// ErrSubjectIsEmpty is returned when an action is given an empty subject.
	// Nothing panics when the action is built; the mistake is reported when the
	// actor performs it.
	ErrSubjectIsEmpty = errors.New("the subject must be provided")

	// ErrQueueIsEmpty is returned when subscribe.To(...).UsingTheQueue is given
	// an empty queue. Nothing panics when the action is built; the mistake is
	// reported when the actor performs it.
	ErrQueueIsEmpty = errors.New("the queue must be provided")
)
