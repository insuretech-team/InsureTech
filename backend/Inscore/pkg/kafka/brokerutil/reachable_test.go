package brokerutil

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReachableBrokers_FiltersUnavailableEntries(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	reachable, unreachable := ReachableBrokers(
		[]string{
			listener.Addr().String(),
			"127.0.0.1:1",
		},
		200*time.Millisecond,
	)

	require.Equal(t, []string{listener.Addr().String()}, reachable)
	require.Equal(t, []string{"127.0.0.1:1"}, unreachable)
}
