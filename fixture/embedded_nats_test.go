package fixture_test

import (
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-nats/v2/fixture"
)

func TestCreateEmbeddedNATSServer(t *testing.T) {
	natsServer := fixture.CreateEmbeddedNATSServer(t)
	assert.NotNil(t, natsServer)

	messageReceived := false
	var mutex sync.Mutex
	sub, err := natsServer.Subscribe("test.subject", func(msg *nats.Msg) {
		assert.Equal(t, "test message", string(msg.Data))
		mutex.Lock()
		defer mutex.Unlock()
		messageReceived = true
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, sub.Unsubscribe())
	}()

	require.NoError(t, natsServer.Publish("test.subject", []byte("test message")))
	assert.Eventually(t, func() bool {
		mutex.Lock()
		defer mutex.Unlock()
		return messageReceived
	}, 2*time.Second, 100*time.Millisecond)
}

func TestRunNATSServer(t *testing.T) {
	natsServer := fixture.RunNATSServer(t)
	assert.NotNil(t, natsServer)
	assert.True(t, natsServer.ReadyForConnections(2*time.Second))
}

func TestRunNATSServer_MultipleInstancesDoNotConflict(t *testing.T) {
	server1 := fixture.RunNATSServer(t)
	server2 := fixture.RunNATSServer(t)
	assert.NotNil(t, server1)
	assert.NotNil(t, server2)
	assert.NotEqual(t, server1.ClientURL(), server2.ClientURL())
}
