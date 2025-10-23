package question_test

import (
	"testing"

	"github.com/grandper/go-screenplay-nats/ability"
	"github.com/grandper/go-screenplay-nats/action"
	"github.com/grandper/go-screenplay-nats/fixture"
	"github.com/grandper/go-screenplay-nats/question"
	action2 "github.com/grandper/go-screenplay/action"
	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/screenplay"
	"github.com/stretchr/testify/assert"
)

func TestBodyOfTheLastResponseQuestion(t *testing.T) {
	nc := fixture.CreateEmbeddedNATSServer(t)
	const subject = "test-subject"
	const content = "Hello World"

	t.Run("returns the body of the last message", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		assert.NoError(t, adam.AttemptsTo(
			action.SubscribeTo(subject),
			action.Publish(ability.Message().To(subject).Containing([]byte(content))),
			action2.PauseFor(1).Second().Because("we wait for the message to be received"),
			action.ReceiveFrom(subject)))
		assert.NoError(t, adam.AttemptsTo(see.The(question.BodyOfTheLastMessage().From(subject), contains.TheText(content))))
	})

	t.Run("fails to get the body of the last message when none has been sent", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateWithNATS().Using(nc))
		assert.Error(t, adam.AttemptsTo(see.The(question.BodyOfTheLastMessage(), contains.TheText(content))))
	})

	t.Run(
		"fails to get the body of the last message if the actor does not have the ability CommunicateWithNATS",
		func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			assert.Error(
				t,
				adam.AttemptsTo(see.The(question.BodyOfTheLastMessage(), contains.TheText(content))),
			)
		},
	)

	t.Run("implements the stringer interface", func(t *testing.T) {
		body := question.BodyOfTheLastMessage()
		assert.Equal(t, "body of the last message", body.String())

		body = question.BodyOfTheLastMessage().From(subject)
		assert.Equal(t, "body of the last message from test-subject", body.String())
	})
}
