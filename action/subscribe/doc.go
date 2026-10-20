// Package subscribe holds the actions that subscribe the actor to a subject:
//
//	err := anActor.AttemptsTo(subscribe.To("orders"))
//	err := anActor.AttemptsTo(subscribe.To("orders").InTheQueueGroup("workers").Asynchronously())
//	err := anActor.AttemptsTo(subscribe.To("orders").AndUnsubscribeAfter(10).Messages())
//	err := anActor.AttemptsTo(subscribe.ToTheSubjectsStartingWith("orders"))
package subscribe
