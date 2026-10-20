// Package subscribe holds the actions that subscribe the actor to a subject:
//
//	err := anActor.AttemptsTo(subscribe.To("orders"))
//	err := anActor.AttemptsTo(subscribe.To("orders").UsingTheQueue("workers").Asynchronously())
//	err := anActor.AttemptsTo(subscribe.To("orders").AndUnsubscribeAfter(10).Messages())
//	err := anActor.AttemptsTo(subscribe.ToTheSubjectStartingWith("orders"))
package subscribe
