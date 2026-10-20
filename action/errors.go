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

	// ErrQueueIsEmpty is returned when subscribe.To(...).InTheQueueGroup is given
	// an empty queue. Nothing panics when the action is built; the mistake is
	// reported when the actor performs it.
	ErrQueueIsEmpty = errors.New("the queue must be provided")

	// ErrInvalidHeaders is wrapped by the error WithTheHeader and WithTheHeaders
	// return when they are given headers they cannot use: an empty key, no header
	// at all, or a key without its value. Nothing panics when the action is
	// built; the mistake is reported when the actor performs it.
	ErrInvalidHeaders = errors.New("invalid headers")
)
