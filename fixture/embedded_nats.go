package fixture

import (
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

const connectionTimeout = 2 * time.Second

// ConnectToAnEmbeddedNATSServer runs an embedded NATS server, shut down when
// the test ends, and returns a connection to it.
func ConnectToAnEmbeddedNATSServer(t *testing.T) *nats.Conn {
	t.Helper()

	natsServer := RunAnEmbeddedNATSServer(t)

	natsConnection, err := nats.Connect(natsServer.ClientURL())
	require.NoError(t, err)
	return natsConnection
}

// RunAnEmbeddedNATSServer runs an embedded NATS server on a free port, shut
// down when the test ends, and returns it.
func RunAnEmbeddedNATSServer(t *testing.T) *server.Server {
	t.Helper()

	opts := natstest.DefaultTestOptions
	opts.Port = -1
	natsServer := natstest.RunServer(&opts)
	t.Cleanup(natsServer.Shutdown)
	require.True(t, natsServer.ReadyForConnections(connectionTimeout))

	return natsServer
}
