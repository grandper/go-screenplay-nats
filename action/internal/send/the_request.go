package send

import (
	"strings"
	"time"

	"github.com/grandper/go-screenplay/v2/action/timing"
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-nats/v2/action/internal/message"
	impl "github.com/grandper/go-screenplay-nats/v2/internal/ability"
	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

// defaultTimeout is how long the actor waits for the reply when TimingOutAfter
// is not written.
const defaultTimeout = 5 * time.Second

// TheRequest sends a request holding the text to the subject given with To
// and waits for its reply, at most five seconds unless TimingOutAfter is
// written; it reads "send the request ping to orders timing out after 2
// seconds". The reply is stored under the subject of the request, where
// messages.FromTheSubject finds it:
//
//	send.TheRequest("ping").To("orders")
//	send.TheRequest("ping").To("orders").WithTheHeader("Key", "value").TimingOutAfter(2).Seconds()
func TheRequest(text string) *TheRequestAction {
	return newRequest("send the request "+verdict.Describe(text), []byte(text))
}

// TheBinaryRequest sends a request holding the bytes to the subject given
// with To and waits for its reply; it reads "send the binary request data to
// orders":
//
//	send.TheBinaryRequest(data).To("orders")
func TheBinaryRequest(data []byte) *TheRequestAction {
	return newRequest("send the binary request "+verdict.Describe(data), data)
}

// AnEmptyRequest sends a request with no payload to the subject given with
// To and waits for its reply; it reads "send an empty request to orders":
//
//	send.AnEmptyRequest().To("orders")
func AnEmptyRequest() *TheRequestAction {
	return newRequest("send an empty request", nil)
}

func newRequest(description string, data []byte) *TheRequestAction {
	return &TheRequestAction{
		description: description,
		message:     message.NewBuilder(data),
		timeout:     defaultTimeout,
	}
}

// TheRequestAction sends a request and waits for its reply.
type TheRequestAction struct {
	description    string
	message        *message.Builder
	timeout        time.Duration
	timeoutBuilder *timing.DurationBuilder[TheRequestAction]
	timeoutWritten bool
}

// To sets the subject the request is sent to; it reads "to orders".
func (a *TheRequestAction) To(subject string) *TheRequestAction {
	a.message.To(subject)

	return a
}

// WithTheHeader adds a header to the request; it reads "with the header
// Key = value".
func (a *TheRequestAction) WithTheHeader(key, value string) *TheRequestAction {
	a.message.WithTheHeader(key, value)

	return a
}

// WithTheHeaders adds headers to the request, written as key, value, key,
// value; it reads "with the headers A = 1, B = 2".
func (a *TheRequestAction) WithTheHeaders(args ...string) *TheRequestAction {
	a.message.WithTheHeaders(args...)

	return a
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
// action.ErrSubjectIsEmpty when no subject was given, and with an error
// wrapping action.ErrInvalidHeaders when the headers were written wrongly.
func (a *TheRequestAction) PerformAs(theActor *screenplay.Actor) error {
	if err := a.message.Err(); err != nil {
		return err
	}

	return impl.TheAbilityOf(theActor).ToSendTheRequest(theActor.Context(), a.message.Msg(), a.timeout)
}

// String describes the action, word for word as it was written.
func (a *TheRequestAction) String() string {
	parts := []string{a.description}
	if chain := a.message.Describe(); chain != "" {
		parts = append(parts, chain)
	}
	if a.timeoutWritten {
		parts = append(parts, "timing out after "+a.timeoutBuilder.String())
	}

	return strings.Join(parts, " ")
}

// TheRequestAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheRequestAction)(nil)
