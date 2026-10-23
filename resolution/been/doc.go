// Package been holds the resolutions about the connection of the actor and
// the messages it received. Each one reads after the verb ShouldHave:
//
//	err := anActor.Should(see.The(connection.ToTheNATSServer()).ShouldHave(been.Closed()))
//	err := anActor.Should(see.The(messages.FromTheSubject("orders")).ShouldHave(been.Received()))
//	err := anActor.Should(see.The(messages.FromAnySubject()).ShouldHave(been.ReceivedFrom("orders")))
//
// They read "see that the connection to the NATS server should have been closed" and "see
// that the NATS messages from the subject orders should have been received".
package been
