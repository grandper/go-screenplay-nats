// Package fixture holds the embedded NATS server that lets a scenario run
// without an external broker:
//
//	nc := fixture.ConnectToAnEmbeddedNATSServer(t)
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().UsingTheConnection(nc))
package fixture
