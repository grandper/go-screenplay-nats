// Package question holds the questions about what the actor received and
// about its connection:
//
//	messages, err := anActor.AsksFor(question.NATSMessages())             // a []*nats.Msg
//	messages, err := anActor.AsksFor(question.NATSMessagesFrom("orders")) // a []*nats.Msg
//	connection, err := anActor.AsksFor(question.NATSConnection())         // an *ability.Connection
//
// A question composes with the core questions and the resolutions of the
// extension:
//
//	err := anActor.Should(see.The(question.NATSMessagesFrom("orders")).ShouldHave(been.Received()))
//	err := anActor.Should(see.The(question.NATSConnection()).ShouldHave(been.Closed()))
package question
