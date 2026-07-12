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

// CreateEmbeddedNATSServer creates an embedded NATS server for testing purposes.
func CreateEmbeddedNATSServer(t *testing.T) *nats.Conn {
	t.Helper()

	natsServer := RunNATSServer(t)

	natsConnection, err := nats.Connect(natsServer.ClientURL())
	require.NoError(t, err)
	return natsConnection
}

// RunNATSServer starts an embedded NATS server for testing purposes.
func RunNATSServer(t *testing.T) *server.Server {
	t.Helper()

	opts := natstest.DefaultTestOptions
	opts.Port = -1
	natsServer := natstest.RunServer(&opts)
	t.Cleanup(natsServer.Shutdown)
	require.True(t, natsServer.ReadyForConnections(connectionTimeout))

	return natsServer
}
