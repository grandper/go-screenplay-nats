// Package ability holds the implementation of the ability to communicate with
// NATS. Only the actions and the questions of the extension use its
// operations, through TheAbilityOf; a test gives the ability to an actor
// through the public package:
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
package ability
