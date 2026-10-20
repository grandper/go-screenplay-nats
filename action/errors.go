package action

import "errors"

var (
	// ErrURLIsEmpty is returned when the URL for the NATS connection is empty.
	ErrURLIsEmpty = errors.New("the URL for the NATS connection cannot be empty")

	// ErrSubjectIsEmpty is returned when the subject is empty.
	ErrSubjectIsEmpty = errors.New("subject is empty")

	// ErrQueueIsEmpty is returned when the queue is empty.
	ErrQueueIsEmpty = errors.New("queue is empty")
)
