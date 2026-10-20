# go-screenplay-nats
This project is an extension of [go-screenplay](https://github.com/grandper/go-screenplay) to drive a NATS server. It bundles the ability to communicate with NATS together with the actions, questions and resolutions that go with it, so you test a messaging scenario with the same screenplay vocabulary you use everywhere else. A test reads as English, and everything the actor does is narrated:
```go
actor := screenplay.ActorNamed("NATS User").WhoCan(ability.CommunicateWithNATS().Using(nc))

err := actor.AttemptsTo(subscribe.To("orders"))
// NATS User attempts to subscribe to "orders"
err = actor.Should(see.The(question.NATSMessagesFrom("orders")).ShouldHave(been.Received()))
// NATS User should see that the NATS messages from "orders" should have been received
```

## Installation
```sh
go get github.com/grandper/go-screenplay-nats/v2
```
It requires go-screenplay v2.

## Packages
The symbols shown below live in these packages:

| Package | Import path | Provides |
|---------|-------------|----------|
| `ability` | `github.com/grandper/go-screenplay-nats/v2/ability` | `CommunicateWithNATS`, `Connection`, `ErrNATSConnectionNotSet`, `ErrFailedToConnect`, `ErrAlreadyConnected`, `ErrSubscriptionAlreadyExists`, `ErrNoSubscriptionFound` |
| `action` | `github.com/grandper/go-screenplay-nats/v2/action` | `ErrURLIsEmpty`, `ErrSubjectIsEmpty`, `ErrQueueIsEmpty` |
| `connect` | `github.com/grandper/go-screenplay-nats/v2/action/connect` | `ToTheNATSServerAt` |
| `subscribe` | `github.com/grandper/go-screenplay-nats/v2/action/subscribe` | `To`, `ToTheSubjectStartingWith` |
| `unsubscribe` | `github.com/grandper/go-screenplay-nats/v2/action/unsubscribe` | `From` |
| `publish` | `github.com/grandper/go-screenplay-nats/v2/action/publish` | `TheMessage` |
| `send` | `github.com/grandper/go-screenplay-nats/v2/action/send` | `TheRequest` |
| `receive` | `github.com/grandper/go-screenplay-nats/v2/action/receive` | `TheNextMessageFrom` |
| `flush` | `github.com/grandper/go-screenplay-nats/v2/action/flush` | `TheNATSConnection` |
| `question` | `github.com/grandper/go-screenplay-nats/v2/question` | `NATSConnection`, `NATSMessages`, `NATSMessagesFrom` |
| `fixture` (of `question`) | `github.com/grandper/go-screenplay-nats/v2/question/fixture` | `CreateFakeConnection`, `CreateFakeClosedConnection`, `CreateFakeDisconnectedConnection`, `CreateFakeReconnectedConnection` |
| `been` | `github.com/grandper/go-screenplay-nats/v2/resolution/been` | `Closed`, `Disconnected`, `Reconnected`, `Received`, `ReceivedFrom` |
| `fixture` | `github.com/grandper/go-screenplay-nats/v2/fixture` | `CreateEmbeddedNATSServer`, `RunNATSServer` |

Assertions are written with the core screenplay packages (`action/see`, `question/last`, `question/number`, `resolution/be`, ...). The extension's `question` and `fixture` packages share their name with core packages: when a file uses both, keep the core name and alias the extension package with the `nats` prefix.

## Ability
This package provides the `CommunicateWithNATS` ability, which allows an actor to interact with a NATS server. The `CommunicateWithNATS` ability is given an existing NATS connection with `Using(nc)`, or the URL of a server with `UsingServerAt("nats://localhost:4222")`.

First, import the package:
```go
import (
    "github.com/nats-io/nats.go"
    "github.com/grandper/go-screenplay/v2/screenplay"
    "github.com/grandper/go-screenplay/v2/action/see"
    "github.com/grandper/go-screenplay-nats/v2/ability"
    "github.com/grandper/go-screenplay-nats/v2/action/connect"
    "github.com/grandper/go-screenplay-nats/v2/action/flush"
    "github.com/grandper/go-screenplay-nats/v2/action/publish"
    "github.com/grandper/go-screenplay-nats/v2/action/receive"
    "github.com/grandper/go-screenplay-nats/v2/action/send"
    "github.com/grandper/go-screenplay-nats/v2/action/subscribe"
    "github.com/grandper/go-screenplay-nats/v2/action/unsubscribe"
    "github.com/grandper/go-screenplay-nats/v2/fixture"
    "github.com/grandper/go-screenplay-nats/v2/question"
    "github.com/grandper/go-screenplay-nats/v2/resolution/been"
)
```
Then create a connection to the NATS server:
```go
nc, _ := nats.Connect(nats.DefaultURL)
nc, _ := nats.Connect("nats://localhost:4222", nats.UserInfo("username", "password"))
```
Finally, create an actor and give it the `CommunicateWithNATS` ability:
```go
actor := screenplay.ActorNamed("NATS User").WhoCan(ability.CommunicateWithNATS().Using(nc))
```
Alternatively, let the ability open the connection itself from a server URL. The connection is opened the first time the actor uses the ability, so an unreachable server is reported by the first action, with an error wrapping `ability.ErrFailedToConnect`:
```go
actor := screenplay.ActorNamed("NATS User").WhoCan(ability.CommunicateWithNATS().UsingServerAt("nats://localhost:4222"))
```
An actor without a connection or a URL fails every action with `ability.ErrNATSConnectionNotSet`, until it performs the `connect.ToTheNATSServerAt` action.

> **_NOTE:_** When the actor exits (`actor.Exit()`), the ability forgets what it holds: every subscription is dropped and the connection is closed, whether it was given with `Using` or opened by the ability itself.

`Using` and `UsingServerAt` are the only methods a test can call on the ability. Its operations live in `internal/ability`, so only an action or a question of this package uses them, in one statement written in `PerformAs` or `AnsweredBy` that reads "the ability of the actor to publish":
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
`subscribe.ToTheSubjectStartingWith` subscribes to every subject starting with a prefix, which is the subject `prefix.>`:
```go
err := actor.AttemptsTo(subscribe.ToTheSubjectStartingWith("orders"))
// NATS User attempts to subscribe to the subject starting with "orders"
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
err := actor.AttemptsTo(subscribe.To("orders").UsingTheQueue("workers"))
// NATS User attempts to subscribe to "orders" using the queue "workers"
```
The subscription can end by itself after a number of messages:
```go
err := actor.AttemptsTo(subscribe.To("orders").AndUnsubscribeAfter(10).Messages())
// NATS User attempts to subscribe to "orders" and unsubscribe after 10 messages
err := actor.AttemptsTo(subscribe.To("orders").AndUnsubscribeAfter(1).Message())
// NATS User attempts to subscribe to "orders" and unsubscribe after 1 message
```
Subscribing twice to the same subject fails with `ability.ErrSubscriptionAlreadyExists`.

### unsubscribe.From
Unsubscribes the actor from a subject it subscribed to:
```go
err := actor.AttemptsTo(unsubscribe.From("orders"))
// NATS User attempts to unsubscribe from "orders"
```
It fails with an error wrapping `ability.ErrNoSubscriptionFound` when the actor did not subscribe to the subject.

### publish.TheMessage
Publishes a message to the subject given with `To`. The payload (`Containing`), the headers (`WithTheHeader`, `WithTheHeaders`) and the subject the message is to be replied to (`ToBeRepliedTo`) are optional, and narrated in the order they are written:
```go
err := actor.AttemptsTo(publish.TheMessage().To("orders"))
// NATS User attempts to publish the message to "orders"
err := actor.AttemptsTo(publish.TheMessage().To("orders").Containing([]byte("Hello, World!")))
// NATS User attempts to publish the message to "orders" containing "Hello, World!"
err := actor.AttemptsTo(publish.TheMessage().
    To("orders").
    Containing([]byte("Hello, World!")).
    WithTheHeader("Key", "value").
    WithTheHeaders(nats.Header{"Other": []string{"1", "2"}}).
    ToBeRepliedTo("inbox"))
// NATS User attempts to publish the message to "orders" containing "Hello, World!" with the header "Key" = "value" with the headers "Other" = "1", "Other" = "2" to be replied to "inbox"
```
`WithTheHeaders` narrates its headers sorted by key.

### send.TheRequest
Sends a request to the subject given with `To` and waits for its reply, which is stored under the subject of the request, where `question.NATSMessagesFrom(subject)` finds it. The payload and the headers are given as for `publish.TheMessage`. The actor waits at most five seconds, unless `TimingOutAfter` is written:
```go
err := actor.AttemptsTo(send.TheRequest().To("orders").Containing([]byte("ping")))
// NATS User attempts to send the request to "orders" containing "ping"
err := actor.AttemptsTo(send.TheRequest().To("orders").TimingOutAfter(2).Seconds())
// NATS User attempts to send the request to "orders" timing out after 2 seconds
err := actor.AttemptsTo(send.TheRequest().To("orders").TimingOutAfter(500).Milliseconds())
// NATS User attempts to send the request to "orders" timing out after 500 milliseconds
```
A request nobody answers fails with `nats.ErrNoResponders`; a reply that does not come in time fails with an error wrapping `nats.ErrTimeout`. A deadline or a cancellation set on the actor's context stops the request as well.

### receive.TheNextMessageFrom
Receives the next message of a synchronous subscription and stores it, where `question.NATSMessagesFrom(subject)` finds it. The actor waits at most its timeout (two seconds by default, see `screenplay.WithTimeout`):
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

### NATSConnection
`question.NATSConnection()` answers with an `*ability.Connection`, a snapshot of what happened to the connection of the actor, read through `WasClosed()`, `WasDisconnected()` and `WasReconnected()`. It is taken when the question is answered and does not change afterwards, so ask the question again to observe a later state:
```go
err := actor.Should(see.The(question.NATSConnection()).ShouldHave(been.Closed()))
// NATS User should see that the NATS connection should have been closed
err := actor.Should(see.The(question.NATSConnection()).ShouldHave(been.Disconnected()))
err := actor.Should(see.The(question.NATSConnection()).ShouldHave(been.Reconnected()))
connection, err := actor.AsksFor(question.NATSConnection()) // connection is an *ability.Connection
```

### NATSMessages
`question.NATSMessages()` answers with the `[]*nats.Msg` received by the actor, in the order they were received: the messages of its subscriptions and the replies to its requests. `question.NATSMessagesFrom(subject)` keeps the messages of the subscription to that subject, or the replies to the requests sent to it:
```go
err := actor.Should(see.The(question.NATSMessages()).ShouldHave(been.Received()))
// NATS User should see that the NATS messages should have been received
err := actor.Should(see.The(question.NATSMessagesFrom("orders")).ShouldHave(been.Received()))
// NATS User should see that the NATS messages from "orders" should have been received
err := actor.Should(see.The(number.Of(question.NATSMessagesFrom("orders"))).Should(be.EqualTo(2)))
lastOrder, err := actor.AsksFor(last.Of(question.NATSMessagesFrom("orders"))) // lastOrder is a *nats.Msg
```

### Fake answers
The `question/fixture` package builds the answer of `NATSConnection` without a server, to test code that receives it:
```go
connection := fixture.CreateFakeConnection()             // nothing happened to it
connection := fixture.CreateFakeClosedConnection()       // disconnected, then closed
connection := fixture.CreateFakeDisconnectedConnection() // disconnected, not reconnected yet
connection := fixture.CreateFakeReconnectedConnection()  // disconnected, then reconnected
```

## Resolutions
The `been` package holds the resolutions of this extension. Each one reads after the verb `ShouldHave`, and a negation lives only in the verb: `ShouldNotHave(been.Closed())` reads "should not have been closed".

`been.Closed()`, `been.Disconnected()` and `been.Reconnected()` inspect the answer of `question.NATSConnection()`. A connection that has been closed or reconnected has been disconnected first:
```go
err := actor.Should(see.The(question.NATSConnection()).ShouldHave(been.Closed()))
// NATS User should see that the NATS connection should have been closed
// fails: NATS User failed to see that the NATS connection should have been closed: it has not been closed
err := actor.Should(see.The(question.NATSConnection()).ShouldNotHave(been.Disconnected()))
// NATS User should see that the NATS connection should not have been disconnected
```
`been.Received()` and `been.ReceivedFrom(subject)` inspect the answer of `question.NATSMessages()` or `question.NATSMessagesFrom(subject)`. The subject of `ReceivedFrom` is the one a message was published to, so a message received through a wildcard subscription matches its own subject:
```go
err := actor.Should(see.The(question.NATSMessagesFrom("orders")).ShouldHave(been.Received()))
// NATS User should see that the NATS messages from "orders" should have been received
// fails: NATS User failed to see that the NATS messages from "orders" should have been received: no message has been received
err := actor.Should(see.The(question.NATSMessages()).ShouldHave(been.ReceivedFrom("orders.created")))
// NATS User should see that the NATS messages should have been received from "orders.created"
```

## Testing with an embedded NATS server
The `fixture` package spins up an in-process NATS server so scenarios can run without an external broker. Both helpers take a `*testing.T` and register their own cleanup, so the server is shut down automatically when the test ends:
```go
// Returns a ready-to-use *nats.Conn connected to the embedded server.
nc := fixture.CreateEmbeddedNATSServer(t)

// Returns the underlying *server.Server if you need to connect several actors
// or inspect the server directly.
srv := fixture.RunNATSServer(t)
```

## Putting it all together
A complete scenario, from spinning up the embedded server to asserting on received messages:
```go
func TestReceivesAPublishedMessage(t *testing.T) {
    nc := fixture.CreateEmbeddedNATSServer(t)
    actor := screenplay.ActorNamed("NATS User").WhoCan(ability.CommunicateWithNATS().Using(nc))

    require.NoError(t, actor.AttemptsTo(
        subscribe.To("greetings"),
        flush.TheNATSConnection(), // ensure the subscription reaches the server before publishing
        publish.TheMessage().To("greetings").Containing([]byte("Hello, World!")),
        receive.TheNextMessageFrom("greetings"),
    ))

    require.NoError(t, actor.Should(
        see.The(question.NATSMessagesFrom("greetings")).ShouldHave(been.Received()),
    ))
}
```

## References
- [NATS documentation](https://docs.nats.io)
- [NATS Go client](https://github.com/nats-io/nats.go)
