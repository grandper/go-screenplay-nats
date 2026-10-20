package verdict_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay-nats/v2/internal/verdict"
)

func TestDescribe(t *testing.T) {
	t.Run("quotes a text", func(t *testing.T) {
		assert.Equal(t, `"orders"`, verdict.Describe("orders"))
	})

	t.Run("quotes bytes", func(t *testing.T) {
		assert.Equal(t, `"Hello"`, verdict.Describe([]byte("Hello")))
	})

	t.Run("writes any other value in its default format", func(t *testing.T) {
		assert.Equal(t, "3", verdict.Describe(3))
		assert.Equal(t, "<nil>", verdict.Describe(nil))
	})
}
