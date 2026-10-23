// Package messages holds the questions about the messages the actor received:
//
//	err := anActor.Should(see.The(messages.FromTheSubject("orders")).ShouldHave(been.Received()))
//	err := anActor.Should(see.The(messages.FromAnySubject()).ShouldHave(been.ReceivedFrom("orders.created")))
//
// A question answers with a []*nats.Msg, so it composes with the core questions:
//
//	err := anActor.Should(see.The(number.Of(messages.FromTheSubject("orders"))).Should(be.EqualTo(2)))
//	lastOrder, err := anActor.AsksFor(last.Of(messages.FromTheSubject("orders")))
package messages
