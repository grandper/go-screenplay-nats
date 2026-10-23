// Package ability holds what a test names about the ability to communicate
// with NATS: the answer type of the connection question and the errors of the
// ability. The ability itself is given to an actor through the communicate
// package:
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
//	status, err := anActor.AsksFor(connection.ToTheNATSServer()) // an *ability.Connection
//	errors.Is(err, ability.ErrFailedToConnect)
package ability
