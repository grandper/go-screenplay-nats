// Package ability holds the ability to communicate with NATS:
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
//
// The ability is given an existing connection with Using, or the URL of a
// server with UsingServerAt, reading "communicate with NATS using server at
// nats://localhost:4222":
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(
//		ability.CommunicateWithNATS().UsingServerAt("nats://localhost:4222"),
//	)
package ability
