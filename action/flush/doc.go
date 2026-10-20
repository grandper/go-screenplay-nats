// Package flush holds the action that flushes the buffered publish operations
// of the NATS connection, so that a subscription or a message reaches the
// server before the next action:
//
//	err := anActor.AttemptsTo(flush.TheNATSConnection())
package flush
