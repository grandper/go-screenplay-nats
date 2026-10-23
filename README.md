# go-screenplay-nats

[![Test](https://github.com/grandper/go-screenplay-nats/actions/workflows/go-test.yml/badge.svg)](https://github.com/grandper/go-screenplay-nats/actions/workflows/go-test.yml)
[![Lint](https://github.com/grandper/go-screenplay-nats/actions/workflows/go-lint.yml/badge.svg)](https://github.com/grandper/go-screenplay-nats/actions/workflows/go-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/grandper/go-screenplay-nats/v2.svg)](https://pkg.go.dev/github.com/grandper/go-screenplay-nats/v2)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

This project is an extension of [go-screenplay](https://github.com/grandper/go-screenplay) to drive a NATS server. It bundles the ability to communicate with NATS together with the actions, questions and resolutions that go with it, so you test a messaging scenario with the same screenplay vocabulary you use everywhere else. A test reads as English, and everything the actor does is narrated:
```go
actor := screenplay.ActorNamed("NATS User").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

err := actor.AttemptsTo(subscribe.To("orders"))
// NATS User attempts to subscribe to "orders"
err = actor.Should(see.The(messages.FromTheSubject("orders")).ShouldHave(been.Received()))
// NATS User should see that the NATS messages from the subject "orders" should have been received
```

## Installation
```sh
go get github.com/grandper/go-screenplay-nats/v2
```
It requires go-screenplay v2.

## Ability
This package provides one ability, `communicate.WithNATS()`, which allows an actor to interact with a NATS server and reads "who can communicate with NATS". It is given an existing NATS connection with `UsingTheConnection(nc)`, reading "who can communicate with NATS using the connection nc", or the URL of a server with `UsingTheServerAt("nats://localhost:4222")`, reading "who can communicate with NATS using the server at nats://localhost:4222".

First, import the package:
```go
import (
    "github.com/nats-io/nats.go"
    "github.com/grandper/go-screenplay/v2/screenplay"
    "github.com/grandper/go-screenplay/v2/action/see"
    "github.com/grandper/go-screenplay-nats/v2/ability"
    "github.com/grandper/go-screenplay-nats/v2/ability/communicate"
    "github.com/grandper/go-screenplay-nats/v2/action/connect"
    "github.com/grandper/go-screenplay-nats/v2/action/flush"
    "github.com/grandper/go-screenplay-nats/v2/action/publish"
    "github.com/grandper/go-screenplay-nats/v2/action/receive"
    "github.com/grandper/go-screenplay-nats/v2/action/send"
    "github.com/grandper/go-screenplay-nats/v2/action/subscribe"
    "github.com/grandper/go-screenplay-nats/v2/action/unsubscribe"
    "github.com/grandper/go-screenplay-nats/v2/fixture"
    "github.com/grandper/go-screenplay-nats/v2/question/connection"
    "github.com/grandper/go-screenplay-nats/v2/question/messages"
    "github.com/grandper/go-screenplay-nats/v2/resolution/been"
)
```
Then create a connection to the NATS server:
```go
nc, _ := nats.Connect(nats.DefaultURL)
nc, _ := nats.Connect("nats://localhost:4222", nats.UserInfo("username", "password"))
```
Finally, create an actor who can communicate with NATS:
```go
actor := screenplay.ActorNamed("NATS User").WhoCan(communicate.WithNATS().UsingTheConnection(nc))
```
Alternatively, let the ability open the connection itself from a server URL. The connection is opened the first time the actor uses the ability, so an unreachable server is reported by the first action, with an error wrapping `ability.ErrFailedToConnect`:
```go
actor := screenplay.ActorNamed("NATS User").WhoCan(communicate.WithNATS().UsingTheServerAt("nats://localhost:4222"))
```
An actor without a connection or a URL fails every action with `ability.ErrNATSConnectionNotSet`, "no connection to NATS was given", until it performs the `connect.ToTheNATSServerAt` action.

> **_NOTE:_** When the actor exits (`actor.Exit()`), the ability forgets what it holds: every subscription is dropped and the connection is closed, whether it was given with `UsingTheConnection` or opened by the ability itself.

`UsingTheConnection` and `UsingTheServerAt` are the only methods a test can call on the ability. Its operations live in `internal/ability`, so only an action or a question of this package uses them, in one statement written in `PerformAs` or `AnsweredBy` that reads "the ability of the actor to publish":
```go
err := impl.TheAbilityOf(theActor).ToPublish(message)
```
`TheAbilityOf` never fails: when the actor cannot communicate with NATS, the operation reports `screenplay.ErrActorMissingAbility`.

You can of course use all the available options when creating the NATS connection, see [NATS Go client](https://github.com/nats-io/nats.go) for more details.

## Actions
The actions live in packages named after their verb, so that a statement reads as English: `subscribe.To("orders")` reads "subscribe to orders". Every action is narrated word for word as it was written, after the name of the actor:
```go
err := actor.AttemptsTo(subscribe.To("orders"))
// NATS User attempts to subscribe to "orders"
```
Nothing panics when an action is built. A mistake, such as an empty subject, is reported when the actor performs the action, with an error of the `action` package: `action.ErrSubjectIsEmpty`, `action.ErrQueueIsEmpty` or `action.ErrURLIsEmpty`.

### connect.ToTheNATSServerAt
Connects the actor to a NATS server, when its ability was given neither a connection nor a URL:
```go
err := actor.AttemptsTo(connect.ToTheNATSServerAt("nats://localhost:4222"))
// NATS User attempts to connect to the NATS server at "nats://localhost:4222"
```
It fails with `action.ErrURLIsEmpty` when the URL is empty, with an error wrapping `ability.ErrFailedToConnect` when the server cannot be reached, and with `ability.ErrAlreadyConnected` when the actor already holds a connection.

### subscribe.To
Subscribes the actor to a subject. Wildcards match any token at one level (`*`) or at every following level (`>`):
```go
err := actor.AttemptsTo(subscribe.To("orders"))
// NATS User attempts to subscribe to "orders"
err := actor.AttemptsTo(subscribe.To("orders.*.created"))
err := actor.AttemptsTo(subscribe.To("orders.>"))
```
`subscribe.ToTheSubjectsStartingWith` subscribes to every subject starting with a prefix, which is the subject `prefix.>`:
```go
err := actor.AttemptsTo(subscribe.ToTheSubjectsStartingWith("orders"))
// NATS User attempts to subscribe to the subjects starting with "orders"
```
A subscription is synchronous by default: its messages are received one at a time with `receive.TheNextMessageFrom`. An asynchronous subscription stores its messages as they arrive:
```go
err := actor.AttemptsTo(subscribe.To("orders").Synchronously())
// NATS User attempts to subscribe to "orders" synchronously
err := actor.AttemptsTo(subscribe.To("orders").Asynchronously())
// NATS User attempts to subscribe to "orders" asynchronously
```
A queue group distributes the messages of a subject among its subscribers, each message going to one of them:
```go
err := actor.AttemptsTo(subscribe.To("orders").InTheQueueGroup("workers"))
// NATS User attempts to subscribe to "orders" in the queue group "workers"
```
The subscription can end by itself after a number of messages:
```go
err := actor.AttemptsTo(subscribe.To("orders").AndUnsubscribeAfter(10).Messages())
// NATS User attempts to subscribe to "orders" and unsubscribe after 10 messages
err := actor.AttemptsTo(subscribe.To("orders").AndUnsubscribeAfter(1).Message())
// NATS User attempts to subscribe to "orders" and unsubscribe after 1 message
```
Subscribing twice to the same subject fails with an error wrapping `ability.ErrSubscriptionAlreadyExists`:
```go
err := actor.AttemptsTo(subscribe.To("orders"), subscribe.To("orders"))
// the actor is already subscribed to the subject "orders"
```

### unsubscribe.From
Unsubscribes the actor from a subject it subscribed to:
```go
err := actor.AttemptsTo(unsubscribe.From("orders"))
// NATS User attempts to unsubscribe from "orders"
```
It fails with an error wrapping `ability.ErrNoSubscriptionFound` when the actor did not subscribe to the subject:
```go
err := actor.AttemptsTo(unsubscribe.From("invoices"))
// the actor is not subscribed to the subject "invoices"
```

### publish.TheMessage
Publishes a message to the subject given with `To`. The message holds a text (`TheMessage`), bytes (`TheBinaryMessage`) or nothing (`AnEmptyMessage`):
```go
err := actor.AttemptsTo(publish.TheMessage("Hello, World!").To("orders"))
// NATS User attempts to publish the message "Hello, World!" to "orders"
err := actor.AttemptsTo(publish.TheBinaryMessage([]byte{0x00, 0x01}).To("orders"))
// NATS User attempts to publish the binary message "\x00\x01" to "orders"
err := actor.AttemptsTo(publish.AnEmptyMessage().To("orders"))
// NATS User attempts to publish an empty message to "orders"
```
The headers (`WithTheHeader`, `WithTheHeaders`, written as `key, value, key, value, ...`) and the subject a reply is expected on (`ExpectingAReplyOn`) are optional, and narrated in the order they are written:
```go
err := actor.AttemptsTo(publish.TheMessage("Hello, World!").
    To("orders").
    WithTheHeader("Key", "value").
    WithTheHeaders("Other", "1", "Other", "2").
    ExpectingAReplyOn("inbox"))
// NATS User attempts to publish the message "Hello, World!" to "orders" with the header "Key" = "value" with the headers "Other" = "1", "Other" = "2" expecting a reply on "inbox"
```
Nothing panics when the headers are written wrongly. An empty key, no header at all, or a key without its value is reported when the actor performs the action, with an error wrapping `action.ErrInvalidHeaders`:
```go
err := actor.AttemptsTo(publish.TheMessage("Hello, World!").To("orders").WithTheHeaders("Key"))
// invalid headers: a key is missing its value
```

### send.TheRequest
Sends a request to the subject given with `To` and waits for its reply, which is stored under the subject of the request, where `messages.FromTheSubject(subject)` finds it. The request holds a text (`TheRequest`), bytes (`TheBinaryRequest`) or nothing (`AnEmptyRequest`), and carries headers as a published message does. The actor waits at most five seconds, unless `TimingOutAfter` is written:
```go
err := actor.AttemptsTo(send.TheRequest("ping").To("orders"))
// NATS User attempts to send the request "ping" to "orders"
err := actor.AttemptsTo(send.TheBinaryRequest([]byte{0x00, 0x01}).To("orders"))
// NATS User attempts to send the binary request "\x00\x01" to "orders"
err := actor.AttemptsTo(send.AnEmptyRequest().To("orders").TimingOutAfter(2).Seconds())
// NATS User attempts to send an empty request to "orders" timing out after 2 seconds
err := actor.AttemptsTo(send.TheRequest("ping").To("orders").WithTheHeader("Key", "value").TimingOutAfter(500).Milliseconds())
// NATS User attempts to send the request "ping" to "orders" with the header "Key" = "value" timing out after 500 milliseconds
```
A request nobody answers fails with `nats.ErrNoResponders`; a reply that does not come in time fails with an error wrapping `nats.ErrTimeout`. A deadline or a cancellation set on the actor's context stops the request as well.

### receive.TheNextMessageFrom
Receives the next message of a synchronous subscription and stores it, where `messages.FromTheSubject(subject)` finds it. The actor waits at most its timeout (two seconds by default, see `screenplay.WithTimeout`):
```go
err := actor.AttemptsTo(receive.TheNextMessageFrom("orders"))
// NATS User attempts to receive the next message from "orders"
```
It fails with an error wrapping `ability.ErrNoSubscriptionFound` when the actor did not subscribe to the subject, and with `nats.ErrSyncSubRequired` when the subscription is asynchronous.

### flush.TheNATSConnection
Flushes the buffered publish operations of the NATS connection: the actor waits until the server has received everything it published or subscribed to so far. It is the way to make sure a subscription is in place before another actor publishes:
```go
err := actor.AttemptsTo(flush.TheNATSConnection())
// NATS User attempts to flush the NATS connection
```

## Questions
A question is a noun phrase that the actor asks for, and its answer is typed: `see.The`, `AsksFor` and the core questions such as `last.Of` and `number.Of` work on it without any type assertion.

### connection.ToTheNATSServer
`connection.ToTheNATSServer()` answers with an `*ability.Connection`, a snapshot of what happened to the connection of the actor, read through `WasClosed()`, `WasDisconnected()` and `WasReconnected()`. It is taken when the question is answered and does not change afterwards, so ask the question again to observe a later state:
```go
err := actor.Should(see.The(connection.ToTheNATSServer()).ShouldHave(been.Closed()))
// NATS User should see that the connection to the NATS server should have been closed
err := actor.Should(see.The(connection.ToTheNATSServer()).ShouldHave(been.Disconnected()))
err := actor.Should(see.The(connection.ToTheNATSServer()).ShouldHave(been.Reconnected()))
status, err := actor.AsksFor(connection.ToTheNATSServer()) // status is an *ability.Connection
```

### messages.FromTheSubject
`messages.FromTheSubject(subject)` answers with the `[]*nats.Msg` the actor received from that subject, in the order they were received: the messages of its subscription to the subject, or the replies to the requests sent to it. `messages.FromAnySubject()` answers with every message the actor received, whatever the subject:
```go
err := actor.Should(see.The(messages.FromTheSubject("orders")).ShouldHave(been.Received()))
// NATS User should see that the NATS messages from the subject "orders" should have been received
err := actor.Should(see.The(messages.FromAnySubject()).ShouldHave(been.Received()))
// NATS User should see that the NATS messages from any subject should have been received
err := actor.Should(see.The(number.Of(messages.FromTheSubject("orders"))).Should(be.EqualTo(2)))
lastOrder, err := actor.AsksFor(last.Of(messages.FromTheSubject("orders"))) // lastOrder is a *nats.Msg
```

### Fake answers
The `question/fixture` package builds the answer of `connection.ToTheNATSServer()` without a server, to test code that receives it:
```go
status := fixture.CreateFakeConnection()             // nothing happened to it
status := fixture.CreateFakeClosedConnection()       // disconnected, then closed
status := fixture.CreateFakeDisconnectedConnection() // disconnected, not reconnected yet
status := fixture.CreateFakeReconnectedConnection()  // disconnected, then reconnected
```

## Resolutions
The `been` package holds the resolutions of this extension. Each one reads after the verb `ShouldHave`, and a negation lives only in the verb: `ShouldNotHave(been.Closed())` reads "should not have been closed".

`been.Closed()`, `been.Disconnected()` and `been.Reconnected()` inspect the answer of `connection.ToTheNATSServer()`. A connection that has been closed or reconnected has been disconnected first:
```go
err := actor.Should(see.The(connection.ToTheNATSServer()).ShouldHave(been.Closed()))
// NATS User should see that the connection to the NATS server should have been closed
// fails: NATS User failed to see that the connection to the NATS server should have been closed: it has not been closed
err := actor.Should(see.The(connection.ToTheNATSServer()).ShouldNotHave(been.Disconnected()))
// NATS User should see that the connection to the NATS server should not have been disconnected
```
`been.Received()` and `been.ReceivedFrom(subject)` inspect the answer of `messages.FromAnySubject()` or `messages.FromTheSubject(subject)`. The subject of `ReceivedFrom` is the one a message was published to, so a message received through a wildcard subscription matches its own subject:
```go
err := actor.Should(see.The(messages.FromTheSubject("orders")).ShouldHave(been.Received()))
// NATS User should see that the NATS messages from the subject "orders" should have been received
// fails: NATS User failed to see that the NATS messages from the subject "orders" should have been received: no message has been received
err := actor.Should(see.The(messages.FromAnySubject()).ShouldHave(been.ReceivedFrom("orders.created")))
// NATS User should see that the NATS messages from any subject should have been received from "orders.created"
```

## Testing with an embedded NATS server
The `fixture` package runs an in-process NATS server so scenarios can run without an external broker. Both helpers take a `*testing.T` and register their own cleanup, so the server is shut down when the test ends:
```go
// Connects to an embedded server and returns the *nats.Conn.
nc := fixture.ConnectToAnEmbeddedNATSServer(t)

// Runs an embedded server and returns the *server.Server, to connect several
// actors to it or to inspect it directly.
natsServer := fixture.RunAnEmbeddedNATSServer(t)
```

## Putting it all together
A complete scenario, from spinning up the embedded server to asserting on received messages:
```go
func TestReceivesAPublishedMessage(t *testing.T) {
    nc := fixture.ConnectToAnEmbeddedNATSServer(t)
    actor := screenplay.ActorNamed("NATS User").WhoCan(communicate.WithNATS().UsingTheConnection(nc))

    require.NoError(t, actor.AttemptsTo(
        subscribe.To("greetings"),
        flush.TheNATSConnection(), // ensure the subscription reaches the server before publishing
        publish.TheMessage("Hello, World!").To("greetings"),
        receive.TheNextMessageFrom("greetings"),
    ))

    require.NoError(t, actor.Should(
        see.The(messages.FromTheSubject("greetings")).ShouldHave(been.Received()),
    ))
}
```

## References
- [NATS documentation](https://docs.nats.io)
- [NATS Go client](https://github.com/nats-io/nats.go)
