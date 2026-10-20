// Package fixture holds the embedded NATS server that lets a scenario run
// without an external broker:
//
//	nc := fixture.CreateEmbeddedNATSServer(t)
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
package fixture
