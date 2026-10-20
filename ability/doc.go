// Package ability holds the ability to communicate with NATS:
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
//
// The ability is given an existing connection with UsingTheConnection, or the
// URL of a server with UsingTheServerAt, reading "communicate with NATS using
// the server at nats://localhost:4222":
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(
//		ability.CommunicateWithNATS().UsingTheServerAt("nats://localhost:4222"),
//	)
package ability
