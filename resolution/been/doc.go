// Package been holds the resolutions about the connection of the actor and
// the messages it received. Each one reads after the verb ShouldHave:
//
//	err := anActor.Should(see.The(question.NATSConnection()).ShouldHave(been.Closed()))
//	err := anActor.Should(see.The(question.NATSMessagesFrom("orders")).ShouldHave(been.Received()))
//	err := anActor.Should(see.The(question.NATSMessages()).ShouldHave(been.ReceivedFrom("orders")))
//
// They read "see that the NATS connection should have been closed" and "see
// that the NATS messages from orders should have been received".
package been
