package send

import (
	"strings"
	"time"

	"github.com/grandper/go-screenplay/v2/action/timing"
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action"
	"github.com/grandper/go-screenplay-nats/v2/action/internal/message"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
)

// defaultTimeout is how long the actor waits for the reply when TimingOutAfter
// is not written.
const defaultTimeout = 5 * time.Second

// TheRequest sends a request to the subject given with To, with the payload
// given with Containing and the headers given with WithTheHeader and
// WithTheHeaders, and waits for its reply, at most five seconds unless
// TimingOutAfter is written. The reply is stored under the subject of the
// request, where the question NATSMessagesFrom finds it:
//
//	send.TheRequest().To("orders").Containing([]byte("ping")).TimingOutAfter(2).Seconds()
func TheRequest() *TheRequestAction {
	a := &TheRequestAction{timeout: defaultTimeout}
	a.Builder = message.NewBuilder(a)

	return a
}

// TheRequestAction sends a request and waits for its reply. Its chain methods
// To, Containing, WithTheHeader and WithTheHeaders come from the message
// builder it embeds.
type TheRequestAction struct {
	*message.Builder[TheRequestAction]

	timeout        time.Duration
	timeoutBuilder *timing.DurationBuilder[TheRequestAction]
	timeoutWritten bool
}

// TimingOutAfter bounds the wait for the reply; it reads "timing out after 2
// seconds". A request that outlasts it fails with an error wrapping
// nats.ErrTimeout.
func (a *TheRequestAction) TimingOutAfter(amount int) *timing.DurationBuilder[TheRequestAction] {
	a.timeoutWritten = true
	a.timeoutBuilder = timing.NewDurationBuilder(a, &a.timeout).For(amount)

	return a.timeoutBuilder
}

// PerformAs performs the action as the provided actor. It fails with
// action.ErrSubjectIsEmpty when no subject was given.
func (a *TheRequestAction) PerformAs(theActor *screenplay.Actor) error {
	if a.Subject() == "" {
		return action.ErrSubjectIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToSendTheRequest(theActor.Context(), a.Msg(), a.timeout)
}

// String describes the action, word for word as it was written.
func (a *TheRequestAction) String() string {
	parts := []string{"send the request"}
	if chain := a.Describe(); chain != "" {
		parts = append(parts, chain)
	}
	if a.timeoutWritten {
		parts = append(parts, "timing out after "+a.timeoutBuilder.String())
	}

	return strings.Join(parts, " ")
}

// TheRequestAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheRequestAction)(nil)
