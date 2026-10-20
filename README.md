# go-screenplay-nats
This project is an extension of [go-screenplay](https://github.com/grandper/go-screenplay) to use NATS.

The `Screenplay` pattern is a design pattern for writing high-quality, maintainable automated tests. It emphasizes the use of actors, abilities, tasks, and questions to create a clear and expressive way to describe test scenarios.

In this document, we describe how to develop NATS scenarios using the Screenplay pattern and this package.

## Installation
```sh
go get github.com/grandper/go-screenplay-nats/v2
```

## Ability
This package provides the `CommunicateWithNATS` ability, which allows an actor to interact with a NATS server. The `CommunicateWithNATS` ability is given an existing NATS connection with `Using(nc)`, or the URL of a server with `UsingServerAt("nats://localhost:4222")`.

First, import the package:
```go
import (
    "github.com/nats-io/nats.go"
    "github.com/grandper/go-screenplay/v2/screenplay"
    "github.com/grandper/go-screenplay/v2/action/see"
    "github.com/grandper/go-screenplay-nats/v2/ability"
    "github.com/grandper/go-screenplay-nats/v2/action"
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
An actor without a connection or a URL fails every action with `ability.ErrNATSConnectionNotSet`, until it performs the `ConnectToNATS` action.

> **_NOTE:_** When the actor exits (`actor.Exit()`), the ability forgets what it holds: every subscription is dropped and the connection is closed, whether it was given with `Using` or opened by the ability itself.

You can of course use all the available options when creating the NATS connection, see [NATS Go client](https://github.com/nats-io/nats.go) for more details.

## Actions

### ConnectToNATS
The `ConnectToNATS` action allows an actor to connect to a NATS server. It can be used as follows:
```go
err := actor.AttemptsTo(action.ConnectToNATS().At("nats://localhost:4222"))
```
It fails with `action.ErrURLIsEmpty` when no URL is given, with an error wrapping `ability.ErrFailedToConnect` when the server cannot be reached, and with `ability.ErrAlreadyConnected` when the actor already holds a connection.

### SubscribeTo
The `SubscribeTo` action allows an actor to subscribe to a NATS subject without providing a message handler. It can be used as follows:
```go
err := actor.AttemptsTo(action.SubscribeTo("foobar"))
```
You can of course use wildcards to match any token at any level of the subject:
```go
err := actor.AttemptsTo(action.SubscribeTo("foo.*.baz"))
err := actor.AttemptsTo(action.SubscribeTo("foo.bar.*"))
```
You can also subscribe to all subjects starting with a given prefix using the `SubscribeToTheSubjectStartingWith` action:
```go
err := actor.AttemptsTo(action.SubscribeToTheSubjectStartingWith("foo"))
```
This is equivalent to using the `>` wildcard to match any number of tokens at the end of the subject:
```go
err := actor.AttemptsTo(action.SubscribeTo("foo.>"))
```

You can also use queue groups to distribute messages among multiple subscribers. To do so, use the `UsingTheQueue` action:
```go
err := actor.AttemptsTo(action.SubscribeTo("subject").UsingTheQueue("queue-group"))
```

You can subscribe to messages synchronously or asynchronously. By default, the subscription is synchronous, meaning that the actor will block until a message is received. To subscribe asynchronously, use the `Asynchronously` option:
```go
err := actor.AttemptsTo(action.SubscribeTo("foobar").Synchronously())
err := actor.AttemptsTo(action.SubscribeTo("foobar").Asynchronously())
```

You can also automatically unsubscribe after receiving a certain number of messages using the `AndUnsubscribeAfter` option:
```go
err := actor.AttemptsTo(action.SubscribeTo("foobar").AndUnsubscribeAfter(10).Messages())
```

### UnsubscribeFrom
The `UnsubscribeFrom` action allows an actor to unsubscribe from a NATS subject. It can be used as follows:
```go
err := actor.AttemptsTo(action.UnsubscribeFrom("foobar"))
```

### Publish
Messages are composed of:
- A subject.
- A payload in the form of a byte array.
- Any number of header fields.
- An optional 'reply' address field.

A message is built with `ability.Message()`:
```go
message := ability.Message().
    To("subject").
    Containing([]byte("Hello, World!")).
    WithHeader("key", "value").
    ToBeRepliedTo("reply-subject")
```
The subject can equivalently be set with `From("subject")` instead of `To("subject")`, and additional headers can be added with `WithHeaders(nats.Header{...})`.

The `Publish` action allows an actor to publish a message to a NATS subject. It can be used as follows:
```go
err := actor.AttemptsTo(action.Publish(message))
```

### Request
The `Request` action allows an actor to publish a message and wait for a reply. The response is stored and can later be retrieved with the `NATSMessages` question. By default the request times out after 5 seconds; use `For` to change the timeout:
```go
err := actor.AttemptsTo(action.Request(message))
err := actor.AttemptsTo(action.Request(message).For(2).Seconds())
err := actor.AttemptsTo(action.Request(message).For(500).Milliseconds())
```

### ReceiveFrom
The `ReceiveFrom` action makes an actor receive the next message from a subject it has subscribed to. The received message is stored and can later be retrieved with the `NATSMessages` question. It can be used as follows:
```go
err := actor.AttemptsTo(action.ReceiveFrom("subject"))
```

### Flush
The `Flush` action allows an actor to force a flush of the buffered publish operations on the NATS connection. It can be used as follows:
```go
err := actor.AttemptsTo(action.Flush())
```

## Questions

### NATSConnection
The `NATSConnection` question allows an actor to retrieve the underlying NATS connection. It can be used as follows:
```go
err := actor.Should(see.The(question.NATSConnection()).ShouldHave(been.Closed()))
err := actor.Should(see.The(question.NATSConnection()).ShouldHave(been.Disconnected()))
err := actor.Should(see.The(question.NATSConnection()).ShouldHave(been.Reconnected()))
```

### NATSMessages
The `NATSMessages` question allows an actor to retrieve all messages received on a specific NATS subject. It can be used as follows:
```go
err := actor.Should(see.The(question.NATSMessages().From("subject")).ShouldHave(been.Received()))
err := actor.Should(see.The(question.NATSMessages()).ShouldHave(been.Received().From("subject")))
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
        action.SubscribeTo("greetings"),
        action.Flush(), // ensure the subscription reaches the server before publishing
        action.Publish(ability.Message().To("greetings").Containing([]byte("Hello, World!"))),
        action.ReceiveFrom("greetings"),
    ))

    require.NoError(t, actor.Should(
        see.The(question.NATSMessages().From("greetings")).ShouldHave(been.Received()),
    ))
}
```

# References
- [NATS documentation](https://docs.nats.io)
- [NATS Go client](https://github.com/nats-io/nats.go)
