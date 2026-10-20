// Package action holds the errors of the actions that drive a NATS server. The
// actions live in packages named after their verb:
//
//	connect.ToTheNATSServerAt("nats://localhost:4222")   action/connect
//	subscribe.To("orders")                               action/subscribe
//	unsubscribe.From("orders")                           action/unsubscribe
//	publish.TheMessage().To("orders")                    action/publish
//	send.TheRequest().To("orders")                       action/send
//	receive.TheNextMessageFrom("orders")                 action/receive
//	flush.TheNATSConnection()                            action/flush
//
// For example:
//
//	err := anActor.AttemptsTo(
//		subscribe.To("orders"),
//		publish.TheMessage().To("orders").Containing([]byte("Hello")),
//		receive.TheNextMessageFrom("orders"),
//	)
package action
