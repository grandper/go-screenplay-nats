// Package connection holds the question about the connection of the actor to
// the NATS server:
//
//	err := anActor.Should(see.The(connection.ToTheNATSServer()).ShouldHave(been.Closed()))
//	status, err := anActor.AsksFor(connection.ToTheNATSServer()) // an *ability.Connection
package connection
