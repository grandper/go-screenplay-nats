# go-screenplay-nats
This project is an extension of [go-screenplay](https://github.com/grandper/go-screenplay) to use NATS.

The `Screenplay` pattern is a design pattern for writing high-quality, maintainable automated tests. It emphasizes the use of actors, abilities, tasks, and questions to create a clear and expressive way to describe test scenarios.

In this document, we describe how to develop NATS scenarios using the Screenplay pattern and this package.

## Ability
This package provides the `CommunicateWithNATS` ability, which allows an actor to interact with a NATS server. The `CommunicateWithNATS` ability can be instantiated with a NATS server URL and optional connection options.

First, import the package:
```go
import (
    "github.com/nats-io/nats.go"
    "github.com/grandper/go-screenplay"
    "github.com/grandper/go-screenplay-nats/ability"
    "github.com/grandper/go-screenplay-nats/action"
    "github.com/grandper/go-screenplay-nats/question"
)
```
Then create a connection to the NATS server:
```go
nc, _ := nats.Connect(nats.DefaultURL)
nc, _ := nats.Connect("nats://localhost:4222", nats.UserInfo("username", "password"))
```
Finally, create an actor and give it the `CommunicateWithNATS` ability:
```go
actor := screenplay.AnActorNamed("NATS User").WhoCan(ability.CommunicateWithNATS().Using(nc))
```

> **_NOTE:_** Since the connection is not created by the actor, you are responsible for cleaning up the connection when it is no longer needed.
> ```go
> // Drain connection (Preferred for responders)
> // Close() not needed if this is called.
> nc.Drain()
>
> // Close connection
> nc.Close()
>```

You can of course use all the available options when creating the NATS connection, see [NATS Go client](https://github.com/nats-io/nats.go) for more details.

## Actions

### ConnectToNATS
The `ConnectToNATS` action allows an actor to connect to a NATS server. It can be used as follows:
```go
err := actor.AttemptsTo(action.ConnectToNATS().At("nats://localhost:4222"))
```

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

### SendMessage
Messages are composed of:
- A subject.
- A payload in the form of a byte array.
- Any number of header fields.
- An optional 'reply' address field.

The `SendMessage` action allows an actor to publish a message to a NATS subject. It can be used as follows:
```go
err := actor.AttemptsTo(action.SendMessage([]byte("Hello, World!")).To("subject").WithHeaders(map[string]string{"key": "value"}).WithReply("reply-subject"))
```

### Flush
The `Flush` action allows an actor to force a flush of the buffered publish operations on the NATS connection. It can be used as follows:
```go
err := actor.AttemptsTo(action.Flush())
```

## Questions

### TheBodyOfTheLastMessage
The `TheBodyOfTheLastMessage` question allows an actor to retrieve the body of the last message received on a specific NATS subject. It can be used as follows:
```go
err := actor.Should(see.The(question.BodyOfTheLastMessage().From("subject"), contains.TheText("Hello, World!")))
```

### NATSConnection
The `NATSConnection` question allows an actor to retrieve the underlying NATS connection. It can be used as follows:
```go
err := actor.Should(see.The(question.NATSConnection(), was.Closed()))
err := actor.Should(see.The(question.NATSConnection(), was.Disconnected()))
err := actor.Should(see.The(question.NATSConnection(), was.Reconnected()))
```

### NATSMessages
The `NATSMessages` question allows an actor to retrieve all messages received on a specific NATS subject. It can be used as follows:
```go
err := actor.Should(see.The(question.NATSMessages().From("subject"), were.Received()))
err := actor.Should(see.The(question.NATSMessages(), were.Received().From("subject")))
```

# References
- [NATS documentation](https://docs.nats.io)
- [NATS Go client](https://github.com/nats-io/nats.go)
