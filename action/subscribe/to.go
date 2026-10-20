package subscribe

import impl "github.com/grandper/go-screenplay-nats/v2/action/internal/subscribe"

// To subscribes the actor to the subject, synchronously unless Asynchronously
// is written: the messages of a synchronous subscription are received with
// receive.TheNextMessageFrom, those of an asynchronous one are stored as
// they arrive. Wildcards are allowed in the subject:
//
//	subscribe.To("orders.*.created")
//	subscribe.To("orders").UsingTheQueue("workers").Asynchronously()
//	subscribe.To("orders").AndUnsubscribeAfter(10).Messages()
func To(subject string) *impl.ToAction {
	return impl.To(subject)
}

// ToTheSubjectStartingWith subscribes the actor to every subject starting
// with the prefix, which is the subject "prefix.>":
//
//	subscribe.ToTheSubjectStartingWith("orders")
func ToTheSubjectStartingWith(prefix string) *impl.ToAction {
	return impl.ToTheSubjectStartingWith(prefix)
}
