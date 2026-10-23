// Package communicate holds the ability to communicate with NATS, reading
// "who can communicate with NATS using the connection nc":
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
//
// The ability is given an existing connection with UsingTheConnection, or the
// URL of a server with UsingTheServerAt, reading "who can communicate with
// NATS using the server at nats://localhost:4222":
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(
//		communicate.WithNATS().UsingTheServerAt("nats://localhost:4222"),
//	)
package communicate
